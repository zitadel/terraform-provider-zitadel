package application_v2

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	apppb "github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/application/v2"
	"google.golang.org/protobuf/types/known/durationpb"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
)

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

	// Reject mid-life application_type changes (e.g. oidc -> saml). The
	// Zitadel API does not support converting an existing application from
	// one type to another, so without this check the user would get a
	// confusing wire-level error at apply time. They need to recreate the
	// resource instead.
	if oldType, newType := activeAppType(d); oldType != "" && newType != "" && oldType != newType {
		return diag.Errorf("changing application_type from %q to %q is not supported by the Zitadel API; remove the resource from configuration and add it back to recreate it with the new type", oldType, newType)
	}

	req := &apppb.UpdateApplicationRequest{
		ApplicationId: d.Id(),
		ProjectId:     d.Get(ProjectIDVar).(string),
		Name:          d.Get(NameVar).(string),
	}

	switch {
	case nestedBlock(d, oidcBlockVar) != nil:
		cfg, derr := buildUpdateOIDC(nestedBlock(d, oidcBlockVar))
		if derr != nil {
			return derr
		}
		req.ApplicationType = &apppb.UpdateApplicationRequest_OidcConfiguration{OidcConfiguration: cfg}
	case nestedBlock(d, samlBlockVar) != nil:
		req.ApplicationType = &apppb.UpdateApplicationRequest_SamlConfiguration{
			SamlConfiguration: buildUpdateSAML(nestedBlock(d, samlBlockVar)),
		}
	case nestedBlock(d, apiBlockVar) != nil:
		req.ApplicationType = &apppb.UpdateApplicationRequest_ApiConfiguration{
			ApiConfiguration: buildUpdateAPI(nestedBlock(d, apiBlockVar)),
		}
	default:
		return diag.Errorf("exactly one of oidc, saml, api must be set")
	}

	if _, err := client.UpdateApplication(ctx, req); err != nil {
		return diag.Errorf("failed to update application: %v", err)
	}
	return read(ctx, d, m)
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
		ApplicationType:          apppb.OIDCApplicationType(apppb.OIDCApplicationType_value[cfg[appTypeVar].(string)]),
		AuthMethodType:           apppb.OIDCAuthMethodType(apppb.OIDCAuthMethodType_value[cfg[authMethodTypeVar].(string)]),
		PostLogoutRedirectUris:   toStringSlice(cfg[postLogoutRedirectURIsVar]),
		Version:                  apppb.OIDCVersion(apppb.OIDCVersion_value[cfg[versionVar].(string)]),
		DevelopmentMode:          cfg[devModeVar].(bool),
		AccessTokenType:          apppb.OIDCTokenType(apppb.OIDCTokenType_value[cfg[accessTokenTypeVar].(string)]),
		AccessTokenRoleAssertion: cfg[accessTokenRoleAssertionVar].(bool),
		IdTokenRoleAssertion:     cfg[idTokenRoleAssertionVar].(bool),
		IdTokenUserinfoAssertion: cfg[idTokenUserinfoAssertionVar].(bool),
		ClockSkew:                durationpb.New(dur),
		AdditionalOrigins:        toStringSlice(cfg[additionalOriginsVar]),
		SkipNativeAppSuccessPage: cfg[skipNativeAppSuccessPageVar].(bool),
		BackChannelLogoutUri:     cfg[backChannelLogoutURIVar].(string),
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
	appType := apppb.OIDCApplicationType(apppb.OIDCApplicationType_value[cfg[appTypeVar].(string)])
	authType := apppb.OIDCAuthMethodType(apppb.OIDCAuthMethodType_value[cfg[authMethodTypeVar].(string)])
	tokenType := apppb.OIDCTokenType(apppb.OIDCTokenType_value[cfg[accessTokenTypeVar].(string)])
	accessTokenRoleAssertion := cfg[accessTokenRoleAssertionVar].(bool)
	idTokenRoleAssertion := cfg[idTokenRoleAssertionVar].(bool)
	idTokenUserinfoAssertion := cfg[idTokenUserinfoAssertionVar].(bool)
	skipNative := cfg[skipNativeAppSuccessPageVar].(bool)
	devMode := cfg[devModeVar].(bool)

	// Pass BackChannelLogoutUri as a pointer unconditionally, including
	// when it is an empty string. Because back_channel_logout_uri is
	// Optional+Computed, simply removing the field from HCL leaves the
	// stored value in state untouched. To actually clear a previously set
	// URI the practitioner sets the attribute to "" explicitly; that
	// empty string then needs to reach the server, which only happens if
	// we always send the pointer (a nil pointer would be treated as "no
	// change" by the API).
	backCh := cfg[backChannelLogoutURIVar].(string)

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
		loginVersionVar:             flattenLoginVersion(oidc.GetLoginVersion()),
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

func flattenSAML(_ *schema.ResourceData, saml *apppb.SAMLConfiguration) map[string]interface{} {
	out := map[string]interface{}{
		loginVersionVar: flattenLoginVersion(saml.GetLoginVersion()),
	}
	if xml := saml.GetMetadataXml(); len(xml) > 0 {
		out[metadataXMLVar] = string(xml)
	}
	if url := saml.GetMetadataUrl(); url != "" {
		out[metadataURLVar] = url
	}
	return out
}

// ---------------------------------------------------------------------------
// API builders / flatteners
// ---------------------------------------------------------------------------

func buildCreateAPI(cfg map[string]interface{}) *apppb.CreateAPIApplicationRequest {
	return &apppb.CreateAPIApplicationRequest{
		AuthMethodType: apppb.APIAuthMethodType(apppb.APIAuthMethodType_value[cfg[authMethodTypeVar].(string)]),
	}
}

func buildUpdateAPI(cfg map[string]interface{}) *apppb.UpdateAPIApplicationConfigurationRequest {
	return &apppb.UpdateAPIApplicationConfigurationRequest{
		AuthMethodType: apppb.APIAuthMethodType(apppb.APIAuthMethodType_value[cfg[authMethodTypeVar].(string)]),
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
func activeAppType(d *schema.ResourceData) (oldType, newType string) {
	for _, key := range []string{oidcBlockVar, samlBlockVar, apiBlockVar} {
		oldV, newV := d.GetChange(key)
		if listHasContent(oldV) {
			oldType = key
		}
		if listHasContent(newV) {
			newType = key
		}
	}
	return oldType, newType
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
