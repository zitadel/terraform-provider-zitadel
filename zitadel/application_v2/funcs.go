package application_v2

import (
	"context"
	"encoding/csv"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	apppb "github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/application/v2"
	"google.golang.org/protobuf/types/known/durationpb"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
)

// importApplication implements `terraform import` for zitadel_application_v2.
//
// Import ID format: <app_id[:org_id[:client_secret]]> (positional; supply an
// empty org_id segment if you need to pass a secret without an org, e.g.
// "<app_id>::<secret>").
//
// client_secret is accepted because the v2 GetApplication RPC never returns
// it (it is only emitted once, at create) and the attribute is Computed-only,
// so it cannot be re-supplied via configuration. Without this, importing or
// migrating a secret-bearing app would permanently drop the secret from
// state. Because the secret lives in a per-type nested block, we look up the
// application to learn its type and seed the secret into the matching
// oidc/api block; read() then preserves it.
func importApplication(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	// Parse the import id the same way the provider's importWithAttributes
	// helper does: a ':'-separated, csv-quoted string. Using csv (with
	// LazyQuotes) handles segments the test/import helpers wrap in quotes,
	// and the SemicolonPlaceholder restore handles a literal ':' inside a
	// segment (e.g. an org id or client secret).
	reader := csv.NewReader(strings.NewReader(d.Id()))
	reader.Comma = ':'
	reader.LazyQuotes = true
	parts, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to parse import id %q: %w", d.Id(), err)
	}
	for i := range parts {
		parts[i] = strings.ReplaceAll(parts[i], helper.SemicolonPlaceholder, ":")
	}

	// Reject extra segments so a typo fails loudly instead of silently
	// importing with an unintended org or secret.
	if len(parts) > 3 {
		return nil, fmt.Errorf("invalid import id %q: expected at most 3 segments <app_id[:org_id[:client_secret]]>, got %d", d.Id(), len(parts))
	}

	appID := parts[0]
	if appID == "" {
		return nil, fmt.Errorf("import id must start with the application id, got %q", d.Id())
	}
	// We intentionally accept any non-empty id rather than applying the
	// strict helper.ConvertID format check, and let the server validate it
	// on the subsequent GetApplication call. This avoids rejecting an id at
	// import time purely on a format assumption.
	d.SetId(appID)

	if len(parts) >= 2 && parts[1] != "" {
		if err := d.Set(helper.OrgIDVar, parts[1]); err != nil {
			return nil, err
		}
	}

	if len(parts) == 3 && parts[2] != "" {
		secret := parts[2]
		clientinfo, ok := m.(*helper.ClientInfo)
		if !ok {
			return nil, fmt.Errorf("failed to get client")
		}
		client, err := helper.GetAppV2Client(ctx, clientinfo)
		if err != nil {
			return nil, err
		}
		// Determine the application type so the secret lands in the right
		// nested block; read() will then preserve it.
		resp, err := client.GetApplication(helper.CtxWithOrgID(ctx, d), &apppb.GetApplicationRequest{ApplicationId: appID})
		if err != nil {
			return nil, fmt.Errorf("failed to look up application %q while importing its client_secret: %w", appID, err)
		}
		app := resp.GetApplication()
		switch {
		case app.GetOidcConfiguration() != nil:
			if err := d.Set(oidcBlockVar, []interface{}{map[string]interface{}{clientSecretVar: secret}}); err != nil {
				return nil, err
			}
		case app.GetApiConfiguration() != nil:
			if err := d.Set(apiBlockVar, []interface{}{map[string]interface{}{clientSecretVar: secret}}); err != nil {
				return nil, err
			}
		default:
			return nil, fmt.Errorf("a client_secret was provided on import, but application %q is neither an OIDC nor an API application (only those have a client secret)", appID)
		}
	}

	return []*schema.ResourceData{d}, nil
}

// create dispatches on which nested config block is set in HCL and fills
// the matching oneof branch of CreateApplicationRequest.application_type.
func create(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started create")

	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}
	client, err := helper.GetAppV2Client(ctx, clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}
	// Scope all subsequent API calls to the org_id attribute so middleware
	// metadata is set consistently with the rest of the provider.
	ctx = helper.CtxWithOrgID(ctx, d)

	req := &apppb.CreateApplicationRequest{
		ProjectId: d.Get(ProjectIDVar).(string),
		Name:      d.Get(NameVar).(string),
	}

	switch {
	case nestedBlock(d, oidcBlockVar) != nil:
		oidc := nestedBlock(d, oidcBlockVar)
		cfg, derr := buildCreateOIDC(oidc)
		if derr != nil {
			return derr
		}
		req.ApplicationType = &apppb.CreateApplicationRequest_OidcConfiguration{OidcConfiguration: cfg}
	case nestedBlock(d, samlBlockVar) != nil:
		req.ApplicationType = &apppb.CreateApplicationRequest_SamlConfiguration{
			SamlConfiguration: buildCreateSAML(nestedBlock(d, samlBlockVar)),
		}
	case nestedBlock(d, apiBlockVar) != nil:
		req.ApplicationType = &apppb.CreateApplicationRequest_ApiConfiguration{
			ApiConfiguration: buildCreateAPI(nestedBlock(d, apiBlockVar)),
		}
	default:
		return diag.Errorf("exactly one of oidc, saml, api must be set")
	}

	resp, err := client.CreateApplication(ctx, req)
	if err != nil {
		return diag.Errorf("failed to create application: %v", err)
	}

	d.SetId(resp.GetApplicationId())

	// A successful CreateApplication leaves the app in ACTIVE state.
	// Set it explicitly so the Computed `state` attribute is populated
	// without a tail-call to read().
	if err := d.Set(stateVar, apppb.ApplicationState_APPLICATION_STATE_ACTIVE.String()); err != nil {
		return diag.FromErr(err)
	}

	// Persist everything we can derive from the CreateApplication
	// response directly into state, without tail-calling read(). The v2
	// GetApplication endpoint exhibits a short eventual-consistency
	// window after CreateApplication: an immediate Get can return OK
	// with an empty Application payload, which the defensive guard in
	// read() then treats as "deleted" and clears d.Id(). Terraform's
	// SDK consistency check then fires with "Root object was present,
	// but now absent". project_v2/funcs.go avoids this by not
	// tail-calling read, so do the same here. Terraform's refresh on
	// the next plan picks up server-derived fields once the write has
	// settled.
	//
	// We surface any d.Set failure because client_secret is only ever
	// returned by the API at create time; silently dropping it would
	// strand the practitioner with no way to recover it short of
	// rotating it server-side.
	if oidc := resp.GetOidcConfiguration(); oidc != nil {
		fields := map[string]interface{}{
			clientIDVar:      oidc.GetClientId(),
			clientSecretVar:  oidc.GetClientSecret(),
			noneCompliantVar: oidc.GetNonCompliant(),
		}
		problems := make([]interface{}, 0, len(oidc.GetComplianceProblems()))
		for _, p := range oidc.GetComplianceProblems() {
			problems = append(problems, map[string]interface{}{
				complianceKeyVar:     p.GetKey(),
				complianceMessageVar: p.GetLocalizedMessage(),
			})
		}
		fields[complianceProblemsVar] = problems
		if err := writeNested(d, oidcBlockVar, fields); err != nil {
			return diag.Errorf("failed to persist OIDC config in state: %v", err)
		}
	}
	if api := resp.GetApiConfiguration(); api != nil {
		if err := writeNested(d, apiBlockVar, map[string]interface{}{
			clientIDVar:     api.GetClientId(),
			clientSecretVar: api.GetClientSecret(),
		}); err != nil {
			return diag.Errorf("failed to persist API client credentials in state: %v", err)
		}
	}

	return nil
}

func read(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started read")

	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}
	client, err := helper.GetAppV2Client(ctx, clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}
	// Scope all subsequent API calls to the org_id attribute so middleware
	// metadata is set consistently with the rest of the provider.
	ctx = helper.CtxWithOrgID(ctx, d)

	resp, err := client.GetApplication(ctx, &apppb.GetApplicationRequest{
		ApplicationId: helper.GetID(d, AppIDVar),
	})
	if err != nil && helper.IgnoreIfNotFoundError(err) == nil {
		d.SetId("")
		return nil
	}
	if err != nil {
		return diag.Errorf("failed to get application: %v", err)
	}

	// Defensive guard against the server returning a successful response
	// with an empty Application payload (instead of a proper NotFound).
	// Without this, the d.SetId at the end of read() would clear the
	// resource ID and Terraform would surface a confusing
	// "Root object was present, but now absent" consistency error.
	app := resp.GetApplication()
	if app == nil || app.GetApplicationId() == "" {
		d.SetId("")
		return nil
	}

	if err := d.Set(NameVar, app.GetName()); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set(stateVar, app.GetState().String()); err != nil {
		return diag.FromErr(err)
	}
	// Refresh project_id from the server response. This keeps state
	// accurate after import and corrects any drift if the application
	// was somehow re-parented out-of-band. The v2 Application proto
	// does not carry organization_id / resource_owner, so org_id is
	// preserved from configuration/import rather than refreshed.
	if pid := app.GetProjectId(); pid != "" {
		if err := d.Set(ProjectIDVar, pid); err != nil {
			return diag.FromErr(err)
		}
	}

	// Clear non-matching config blocks before populating the active one.
	// Without this an import (or a stale state from before an out-of-band app
	// type change) could leave two blocks populated and violate the
	// ExactlyOneOf constraint, surfacing as "inconsistent result" errors.
	switch {
	case app.GetOidcConfiguration() != nil:
		if err := d.Set(samlBlockVar, nil); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set(apiBlockVar, nil); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set(oidcBlockVar, []interface{}{flattenOIDC(d, app.GetOidcConfiguration())}); err != nil {
			return diag.FromErr(err)
		}
	case app.GetSamlConfiguration() != nil:
		if err := d.Set(oidcBlockVar, nil); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set(apiBlockVar, nil); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set(samlBlockVar, []interface{}{flattenSAML(d, app.GetSamlConfiguration())}); err != nil {
			return diag.FromErr(err)
		}
	case app.GetApiConfiguration() != nil:
		if err := d.Set(oidcBlockVar, nil); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set(samlBlockVar, nil); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set(apiBlockVar, []interface{}{flattenAPI(d, app.GetApiConfiguration())}); err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(app.GetApplicationId())
	return nil
}

func update(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started update")

	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}
	client, err := helper.GetAppV2Client(ctx, clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}
	// Scope all subsequent API calls to the org_id attribute so middleware
	// metadata is set consistently with the rest of the provider.
	ctx = helper.CtxWithOrgID(ctx, d)

	// A change of application type (oidc/saml/api) is handled at plan time by
	// forceNewOnAppTypeChange, which forces a replacement, so update() is only
	// reached when the type is unchanged.

	req := &apppb.UpdateApplicationRequest{
		ApplicationId: d.Id(),
		ProjectId:     d.Get(ProjectIDVar).(string),
	}

	// Both Name and the config oneof are optional on UpdateApplication
	// ("If not set, the name will not be changed"). Only include the
	// parts that actually changed: resending an unchanged config makes
	// Zitadel reject the call with FailedPrecondition "No changes", which
	// would break a name-only update. This mirrors the d.HasChange gating
	// the v1 application_oidc resource does across its two update RPCs.
	nameChanged := d.HasChange(NameVar)
	if nameChanged {
		req.Name = d.Get(NameVar).(string)
	}

	switch {
	case nestedBlock(d, oidcBlockVar) != nil:
		if d.HasChange(oidcBlockVar) {
			cfg, derr := buildUpdateOIDC(nestedBlock(d, oidcBlockVar))
			if derr != nil {
				return derr
			}
			req.ApplicationType = &apppb.UpdateApplicationRequest_OidcConfiguration{OidcConfiguration: cfg}
		}
	case nestedBlock(d, samlBlockVar) != nil:
		if d.HasChange(samlBlockVar) {
			req.ApplicationType = &apppb.UpdateApplicationRequest_SamlConfiguration{
				SamlConfiguration: buildUpdateSAML(nestedBlock(d, samlBlockVar)),
			}
		}
	case nestedBlock(d, apiBlockVar) != nil:
		if d.HasChange(apiBlockVar) {
			req.ApplicationType = &apppb.UpdateApplicationRequest_ApiConfiguration{
				ApiConfiguration: buildUpdateAPI(nestedBlock(d, apiBlockVar)),
			}
		}
	default:
		return diag.Errorf("exactly one of oidc, saml, api must be set")
	}

	// Nothing we send changed; skip the API call entirely to avoid a
	// spurious "No changes" error from Zitadel. We test the change signal
	// (nameChanged) rather than req.Name to decide this, independent of the
	// name's value. State already equals the plan, so there is nothing to
	// refresh.
	if !nameChanged && req.ApplicationType == nil {
		return nil
	}

	if _, err := client.UpdateApplication(ctx, req); err != nil {
		return diag.Errorf("failed to update application: %v", err)
	}
	// Deliberately do not tail-call read() here. The v2 GetApplication
	// endpoint can lag immediately after UpdateApplication (read-after-
	// write), returning the pre-update values and reverting freshly
	// applied attributes in state, which surfaces as a non-empty plan
	// after apply. Terraform already holds the applied config values;
	// leave them in state and let the next refresh reconcile any
	// server-side computed fields. This matches the v1 application_oidc
	// update, which also returns without re-reading.
	return nil
}

func delete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started delete")

	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}
	client, err := helper.GetAppV2Client(ctx, clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}
	// Scope all subsequent API calls to the org_id attribute so middleware
	// metadata is set consistently with the rest of the provider.
	ctx = helper.CtxWithOrgID(ctx, d)
	if _, err := client.DeleteApplication(ctx, &apppb.DeleteApplicationRequest{
		ApplicationId: d.Id(),
		ProjectId:     d.Get(ProjectIDVar).(string),
	}); err != nil {
		return diag.Errorf("failed to delete application: %v", err)
	}
	return nil
}

// ---------------------------------------------------------------------------
// OIDC builders / flatteners
// ---------------------------------------------------------------------------

func buildCreateOIDC(cfg map[string]interface{}) (*apppb.CreateOIDCApplicationRequest, diag.Diagnostics) {
	respTypes := make([]apppb.OIDCResponseType, 0)
	for _, v := range cfg[responseTypesVar].([]interface{}) {
		respTypes = append(respTypes, apppb.OIDCResponseType(apppb.OIDCResponseType_value[v.(string)]))
	}
	grantTypes := make([]apppb.OIDCGrantType, 0)
	for _, v := range cfg[grantTypesVar].([]interface{}) {
		grantTypes = append(grantTypes, apppb.OIDCGrantType(apppb.OIDCGrantType_value[v.(string)]))
	}

	dur, err := time.ParseDuration(stringOrDefault(cfg[clockSkewVar], "0s"))
	if err != nil {
		return nil, diag.FromErr(err)
	}

	return &apppb.CreateOIDCApplicationRequest{
		RedirectUris:             toStringSlice(cfg[redirectURIsVar]),
		ResponseTypes:            respTypes,
		GrantTypes:               grantTypes,
		ApplicationType:          apppb.OIDCApplicationType(apppb.OIDCApplicationType_value[cfgString(cfg, appTypeVar)]),
		AuthMethodType:           apppb.OIDCAuthMethodType(apppb.OIDCAuthMethodType_value[cfgString(cfg, authMethodTypeVar)]),
		PostLogoutRedirectUris:   toStringSlice(cfg[postLogoutRedirectURIsVar]),
		Version:                  apppb.OIDCVersion(apppb.OIDCVersion_value[cfgString(cfg, versionVar)]),
		DevelopmentMode:          cfgBool(cfg, devModeVar),
		AccessTokenType:          apppb.OIDCTokenType(apppb.OIDCTokenType_value[cfgString(cfg, accessTokenTypeVar)]),
		AccessTokenRoleAssertion: cfgBool(cfg, accessTokenRoleAssertionVar),
		IdTokenRoleAssertion:     cfgBool(cfg, idTokenRoleAssertionVar),
		IdTokenUserinfoAssertion: cfgBool(cfg, idTokenUserinfoAssertionVar),
		ClockSkew:                durationpb.New(dur),
		AdditionalOrigins:        toStringSlice(cfg[additionalOriginsVar]),
		SkipNativeAppSuccessPage: cfgBool(cfg, skipNativeAppSuccessPageVar),
		BackChannelLogoutUri:     cfgString(cfg, backChannelLogoutURIVar),
		LoginVersion:             buildLoginVersion(cfg[loginVersionVar]),
	}, nil
}

func buildUpdateOIDC(cfg map[string]interface{}) (*apppb.UpdateOIDCApplicationConfigurationRequest, diag.Diagnostics) {
	respTypes := make([]apppb.OIDCResponseType, 0)
	for _, v := range cfg[responseTypesVar].([]interface{}) {
		respTypes = append(respTypes, apppb.OIDCResponseType(apppb.OIDCResponseType_value[v.(string)]))
	}
	grantTypes := make([]apppb.OIDCGrantType, 0)
	for _, v := range cfg[grantTypesVar].([]interface{}) {
		grantTypes = append(grantTypes, apppb.OIDCGrantType(apppb.OIDCGrantType_value[v.(string)]))
	}
	dur, err := time.ParseDuration(stringOrDefault(cfg[clockSkewVar], "0s"))
	if err != nil {
		return nil, diag.FromErr(err)
	}
	appType := apppb.OIDCApplicationType(apppb.OIDCApplicationType_value[cfgString(cfg, appTypeVar)])
	authType := apppb.OIDCAuthMethodType(apppb.OIDCAuthMethodType_value[cfgString(cfg, authMethodTypeVar)])
	tokenType := apppb.OIDCTokenType(apppb.OIDCTokenType_value[cfgString(cfg, accessTokenTypeVar)])
	accessTokenRoleAssertion := cfgBool(cfg, accessTokenRoleAssertionVar)
	idTokenRoleAssertion := cfgBool(cfg, idTokenRoleAssertionVar)
	idTokenUserinfoAssertion := cfgBool(cfg, idTokenUserinfoAssertionVar)
	skipNative := cfgBool(cfg, skipNativeAppSuccessPageVar)
	devMode := cfgBool(cfg, devModeVar)

	// Pass BackChannelLogoutUri as a pointer unconditionally, including
	// when it is an empty string. Because back_channel_logout_uri is
	// Optional+Computed, simply removing the field from HCL leaves the
	// stored value in state untouched. To actually clear a previously set
	// URI the practitioner sets the attribute to "" explicitly; that
	// empty string then needs to reach the server, which only happens if
	// we always send the pointer (a nil pointer would be treated as "no
	// change" by the API).
	backCh := cfgString(cfg, backChannelLogoutURIVar)

	return &apppb.UpdateOIDCApplicationConfigurationRequest{
		RedirectUris:             toStringSlice(cfg[redirectURIsVar]),
		ResponseTypes:            respTypes,
		GrantTypes:               grantTypes,
		ApplicationType:          &appType,
		AuthMethodType:           &authType,
		PostLogoutRedirectUris:   toStringSlice(cfg[postLogoutRedirectURIsVar]),
		DevelopmentMode:          &devMode,
		AccessTokenType:          &tokenType,
		AccessTokenRoleAssertion: &accessTokenRoleAssertion,
		IdTokenRoleAssertion:     &idTokenRoleAssertion,
		IdTokenUserinfoAssertion: &idTokenUserinfoAssertion,
		AdditionalOrigins:        toStringSlice(cfg[additionalOriginsVar]),
		ClockSkew:                durationpb.New(dur),
		SkipNativeAppSuccessPage: &skipNative,
		BackChannelLogoutUri:     &backCh,
		LoginVersion:             buildLoginVersion(cfg[loginVersionVar]),
	}, nil
}

func flattenOIDC(d *schema.ResourceData, oidc *apppb.OIDCConfiguration) map[string]interface{} {
	grantTypes := make([]string, 0, len(oidc.GetGrantTypes()))
	for _, g := range oidc.GetGrantTypes() {
		grantTypes = append(grantTypes, g.String())
	}
	respTypes := make([]string, 0, len(oidc.GetResponseTypes()))
	for _, r := range oidc.GetResponseTypes() {
		respTypes = append(respTypes, r.String())
	}
	clockSkew := oidc.GetClockSkew().AsDuration().String()
	if clockSkew == "" {
		clockSkew = "0s"
	}

	problems := make([]interface{}, 0, len(oidc.GetComplianceProblems()))
	for _, p := range oidc.GetComplianceProblems() {
		problems = append(problems, map[string]interface{}{
			complianceKeyVar:     p.GetKey(),
			complianceMessageVar: p.GetLocalizedMessage(),
		})
	}

	out := map[string]interface{}{
		redirectURIsVar:             oidc.GetRedirectUris(),
		responseTypesVar:            respTypes,
		grantTypesVar:               grantTypes,
		appTypeVar:                  oidc.GetApplicationType().String(),
		authMethodTypeVar:           oidc.GetAuthMethodType().String(),
		postLogoutRedirectURIsVar:   oidc.GetPostLogoutRedirectUris(),
		versionVar:                  oidc.GetVersion().String(),
		devModeVar:                  oidc.GetDevelopmentMode(),
		accessTokenTypeVar:          oidc.GetAccessTokenType().String(),
		accessTokenRoleAssertionVar: oidc.GetAccessTokenRoleAssertion(),
		idTokenRoleAssertionVar:     oidc.GetIdTokenRoleAssertion(),
		idTokenUserinfoAssertionVar: oidc.GetIdTokenUserinfoAssertion(),
		clockSkewVar:                clockSkew,
		additionalOriginsVar:        oidc.GetAdditionalOrigins(),
		skipNativeAppSuccessPageVar: oidc.GetSkipNativeAppSuccessPage(),
		backChannelLogoutURIVar:     oidc.GetBackChannelLogoutUri(),
		clientIDVar:                 oidc.GetClientId(),
		noneCompliantVar:            oidc.GetNonCompliant(),
		complianceProblemsVar:       problems,
	}

	// Only set login_version when the server actually returned one. It is
	// Optional+Computed, so writing an empty/null value can produce a
	// perpetual diff; omitting it lets Terraform keep the computed value.
	// Mirrors the v1 application_oidc read behaviour.
	if lv := flattenLoginVersion(oidc.GetLoginVersion()); len(lv) > 0 {
		out[loginVersionVar] = lv
	}

	// Preserve client_secret if previously stored (server doesn't return it on Get).
	if prev := nestedBlock(d, oidcBlockVar); prev != nil {
		if cs, ok := prev[clientSecretVar].(string); ok && cs != "" {
			out[clientSecretVar] = cs
		}
	}
	return out
}

// ---------------------------------------------------------------------------
// SAML builders / flatteners
// ---------------------------------------------------------------------------

func buildCreateSAML(cfg map[string]interface{}) *apppb.CreateSAMLApplicationRequest {
	req := &apppb.CreateSAMLApplicationRequest{
		LoginVersion: buildLoginVersion(cfg[loginVersionVar]),
	}
	if xml, ok := cfg[metadataXMLVar].(string); ok && xml != "" {
		req.Metadata = &apppb.CreateSAMLApplicationRequest_MetadataXml{MetadataXml: []byte(xml)}
	} else if url, ok := cfg[metadataURLVar].(string); ok && url != "" {
		req.Metadata = &apppb.CreateSAMLApplicationRequest_MetadataUrl{MetadataUrl: url}
	}
	return req
}

func buildUpdateSAML(cfg map[string]interface{}) *apppb.UpdateSAMLApplicationConfigurationRequest {
	req := &apppb.UpdateSAMLApplicationConfigurationRequest{
		LoginVersion: buildLoginVersion(cfg[loginVersionVar]),
	}
	if xml, ok := cfg[metadataXMLVar].(string); ok && xml != "" {
		req.Metadata = &apppb.UpdateSAMLApplicationConfigurationRequest_MetadataXml{MetadataXml: []byte(xml)}
	} else if url, ok := cfg[metadataURLVar].(string); ok && url != "" {
		req.Metadata = &apppb.UpdateSAMLApplicationConfigurationRequest_MetadataUrl{MetadataUrl: url}
	}
	return req
}

func flattenSAML(d *schema.ResourceData, saml *apppb.SAMLConfiguration) map[string]interface{} {
	out := map[string]interface{}{}

	// Only set login_version when the server returned one (Optional+Computed;
	// writing an empty value can cause a perpetual diff).
	if lv := flattenLoginVersion(saml.GetLoginVersion()); len(lv) > 0 {
		out[loginVersionVar] = lv
	}

	// A metadata_url-backed application has its metadata fetched and stored
	// by ZITADEL, so GetApplication returns BOTH metadata_url (the source)
	// and metadata_xml (the resolved document). Populating both would
	// violate the metadata_xml/metadata_url ExactlyOneOf constraint and
	// produce an inconsistent-result error or a perpetual diff. Populate
	// only one field:
	//   - if prior state shows which one the practitioner configured, keep
	//     that one (so an xml-configured app is not flipped to the url, or
	//     vice versa);
	//   - on a fresh import with no prior state, prefer metadata_url when
	//     present (it is the smaller source of truth and avoids storing the
	//     large resolved XML), falling back to metadata_xml otherwise.
	xmlConfigured, urlConfigured := false, false
	if prev := nestedBlock(d, samlBlockVar); prev != nil {
		if u, ok := prev[metadataURLVar].(string); ok && u != "" {
			urlConfigured = true
		}
		if x, ok := prev[metadataXMLVar].(string); ok && x != "" {
			xmlConfigured = true
		}
	}

	switch {
	case urlConfigured:
		out[metadataURLVar] = saml.GetMetadataUrl()
	case xmlConfigured:
		out[metadataXMLVar] = string(saml.GetMetadataXml())
	case saml.GetMetadataUrl() != "":
		out[metadataURLVar] = saml.GetMetadataUrl()
	case len(saml.GetMetadataXml()) > 0:
		out[metadataXMLVar] = string(saml.GetMetadataXml())
	}
	return out
}

// ---------------------------------------------------------------------------
// API builders / flatteners
// ---------------------------------------------------------------------------

func buildCreateAPI(cfg map[string]interface{}) *apppb.CreateAPIApplicationRequest {
	return &apppb.CreateAPIApplicationRequest{
		AuthMethodType: apppb.APIAuthMethodType(apppb.APIAuthMethodType_value[cfgString(cfg, authMethodTypeVar)]),
	}
}

func buildUpdateAPI(cfg map[string]interface{}) *apppb.UpdateAPIApplicationConfigurationRequest {
	return &apppb.UpdateAPIApplicationConfigurationRequest{
		AuthMethodType: apppb.APIAuthMethodType(apppb.APIAuthMethodType_value[cfgString(cfg, authMethodTypeVar)]),
	}
}

func flattenAPI(d *schema.ResourceData, api *apppb.APIConfiguration) map[string]interface{} {
	out := map[string]interface{}{
		authMethodTypeVar: api.GetAuthMethodType().String(),
		clientIDVar:       api.GetClientId(),
	}
	if prev := nestedBlock(d, apiBlockVar); prev != nil {
		if cs, ok := prev[clientSecretVar].(string); ok && cs != "" {
			out[clientSecretVar] = cs
		}
	}
	return out
}

// ---------------------------------------------------------------------------
// LoginVersion helpers (shared between OIDC and SAML)
// ---------------------------------------------------------------------------

func buildLoginVersion(raw interface{}) *apppb.LoginVersion {
	list, ok := raw.([]interface{})
	if !ok || len(list) == 0 || list[0] == nil {
		return nil
	}
	item := list[0].(map[string]interface{})

	if v, ok := item[loginV1Var]; ok && v.(bool) {
		return &apppb.LoginVersion{Version: &apppb.LoginVersion_LoginV1{LoginV1: &apppb.LoginV1{}}}
	}
	if v, ok := item[loginV2Var]; ok && v != nil {
		v2list, _ := v.([]interface{})
		if len(v2list) > 0 && v2list[0] != nil {
			v2item := v2list[0].(map[string]interface{})
			var base *string
			if s, ok := v2item[baseURIVar].(string); ok && s != "" {
				base = &s
			}
			return &apppb.LoginVersion{Version: &apppb.LoginVersion_LoginV2{LoginV2: &apppb.LoginV2{BaseUri: base}}}
		}
	}
	return nil
}

func flattenLoginVersion(lv *apppb.LoginVersion) []interface{} {
	if lv == nil {
		return nil
	}
	switch lv.GetVersion().(type) {
	case *apppb.LoginVersion_LoginV1:
		return []interface{}{map[string]interface{}{loginV1Var: true}}
	case *apppb.LoginVersion_LoginV2:
		v2 := lv.GetLoginV2()
		entry := map[string]interface{}{}
		if u := v2.GetBaseUri(); u != "" {
			entry[baseURIVar] = u
		}
		return []interface{}{map[string]interface{}{loginV2Var: []interface{}{entry}}}
	}
	return nil
}

// ---------------------------------------------------------------------------
// Generic helpers
// ---------------------------------------------------------------------------

// nestedBlock returns the single map inside a TypeList(MaxItems=1) block, or
// nil if the block is absent / empty.
func nestedBlock(d *schema.ResourceData, key string) map[string]interface{} {
	raw, ok := d.GetOk(key)
	if !ok {
		return nil
	}
	list, ok := raw.([]interface{})
	if !ok || len(list) == 0 || list[0] == nil {
		return nil
	}
	m, ok := list[0].(map[string]interface{})
	if !ok {
		return nil
	}
	return m
}

// writeNested merges the given fields into the existing nested block (or
// creates the block if missing). Used to persist server-generated secrets
// alongside what the user wrote into HCL. The error from d.Set is
// propagated so callers can surface it instead of silently dropping
// credentials that the server only returns once.
func writeNested(d *schema.ResourceData, key string, fields map[string]interface{}) error {
	cur := nestedBlock(d, key)
	if cur == nil {
		cur = map[string]interface{}{}
	}
	for k, v := range fields {
		cur[k] = v
	}
	return d.Set(key, []interface{}{cur})
}

func toStringSlice(in interface{}) []string {
	if in == nil {
		return nil
	}
	list, ok := in.([]interface{})
	if !ok {
		return nil
	}
	out := make([]string, 0, len(list))
	for _, v := range list {
		out = append(out, v.(string))
	}
	return out
}

// activeAppType inspects the prior and proposed state of the three mutually
// exclusive config blocks and returns the active application type for each.
// An empty string means no block was populated on that side. ResourceData's
// GetChange returns (old, new) pairs even during plan/apply, so this works
// both for detecting the active type during an update and for the
// new-resource case (where old is empty).
// forceNewOnAppTypeChange marks the resource for replacement when the active
// configuration block (oidc/saml/api) changes. The Zitadel API cannot convert
// an existing application from one type to another, so the change is surfaced
// as a plan-time replacement rather than an apply-time failure.
func forceNewOnAppTypeChange(_ context.Context, d *schema.ResourceDiff, _ interface{}) error {
	blocks := []string{oidcBlockVar, samlBlockVar, apiBlockVar}
	var oldType, newType string
	for _, key := range blocks {
		oldV, newV := d.GetChange(key)
		if listHasContent(oldV) {
			oldType = key
		}
		if listHasContent(newV) {
			newType = key
		}
	}
	if oldType != "" && newType != "" && oldType != newType {
		for _, key := range blocks {
			if err := d.ForceNew(key); err != nil {
				return err
			}
		}
	}
	return nil
}

func listHasContent(v interface{}) bool {
	list, ok := v.([]interface{})
	if !ok {
		return false
	}
	return len(list) > 0 && list[0] != nil
}

func stringOrDefault(v interface{}, def string) string {
	if s, ok := v.(string); ok && s != "" {
		return s
	}
	return def
}

// cfgString and cfgBool read a value from a nested-block map with a safe
// type assertion. SDKv2 normally populates every block attribute with its
// zero value, but using ok-checked casts avoids a panic if a key is ever
// absent (e.g. partial state after import or a future schema change).
func cfgString(m map[string]interface{}, key string) string {
	if v, ok := m[key].(string); ok {
		return v
	}
	return ""
}

func cfgBool(m map[string]interface{}, key string) bool {
	if v, ok := m[key].(bool); ok {
		return v
	}
	return false
}
