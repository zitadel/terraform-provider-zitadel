package application_oidc

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

	_, err = client.DeleteApplication(ctx, &apppb.DeleteApplicationRequest{
		ApplicationId: d.Id(),
		ProjectId:     d.Get(ProjectIDVar).(string),
	})
	if err != nil {
		return diag.Errorf("failed to delete applicationOIDC: %v", err)
	}
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

	projectID := d.Get(ProjectIDVar).(string)

	// Build OIDC config update
	oidcConfig := &apppb.UpdateOIDCApplicationConfigurationRequest{}

	if d.HasChange(NameVar) {
		// Name is on the top-level UpdateApplicationRequest, not in OIDC config
	}
	if d.HasChanges(
		redirectURIsVar,
		responseTypesVar,
		grantTypesVar,
		appTypeVar,
		authMethodTypeVar,
		postLogoutRedirectURIsVar,
		devModeVar,
		accessTokenTypeVar,
		accessTokenRoleAssertionVar,
		idTokenRoleAssertionVar,
		idTokenUserinfoAssertionVar,
		clockSkewVar,
		additionalOriginsVar,
		skipNativeAppSuccessPageVar,
		BackChannelLogoutURIVar,
		LoginVersionVar,
	) {
		respTypes := make([]apppb.OIDCResponseType, 0)
		for _, respType := range d.Get(responseTypesVar).([]interface{}) {
			respTypes = append(respTypes, apppb.OIDCResponseType(apppb.OIDCResponseType_value[respType.(string)]))
		}
		grantTypes := make([]apppb.OIDCGrantType, 0)
		for _, grantType := range d.Get(grantTypesVar).([]interface{}) {
			grantTypes = append(grantTypes, apppb.OIDCGrantType(apppb.OIDCGrantType_value[grantType.(string)]))
		}
		dur, err := time.ParseDuration(d.Get(clockSkewVar).(string))
		if err != nil {
			return diag.FromErr(err)
		}

		appType := apppb.OIDCApplicationType(apppb.OIDCApplicationType_value[d.Get(appTypeVar).(string)])
		authType := apppb.OIDCAuthMethodType(apppb.OIDCAuthMethodType_value[d.Get(authMethodTypeVar).(string)])
		accessTokenType := apppb.OIDCTokenType(apppb.OIDCTokenType_value[d.Get(accessTokenTypeVar).(string)])
		accessTokenRoleAssertion := d.Get(accessTokenRoleAssertionVar).(bool)
		idTokenRoleAssertion := d.Get(idTokenRoleAssertionVar).(bool)
		idTokenUserinfoAssertion := d.Get(idTokenUserinfoAssertionVar).(bool)
		skipNative := d.Get(skipNativeAppSuccessPageVar).(bool)

		oidcConfig = &apppb.UpdateOIDCApplicationConfigurationRequest{
			RedirectUris:             interfaceToStringSlice(d.Get(redirectURIsVar)),
			ResponseTypes:            respTypes,
			GrantTypes:               grantTypes,
			ApplicationType:          &appType,
			AuthMethodType:           &authType,
			PostLogoutRedirectUris:   interfaceToStringSlice(d.Get(postLogoutRedirectURIsVar)),
			AccessTokenType:          &accessTokenType,
			AccessTokenRoleAssertion: &accessTokenRoleAssertion,
			IdTokenRoleAssertion:     &idTokenRoleAssertion,
			IdTokenUserinfoAssertion: &idTokenUserinfoAssertion,
			AdditionalOrigins:        interfaceToStringSlice(d.Get(additionalOriginsVar)),
			ClockSkew:                durationpb.New(dur),
			SkipNativeAppSuccessPage: &skipNative,
			LoginVersion:             getLoginVersion(d),
		}
		val := d.Get(BackChannelLogoutURIVar).(string)
		if val != "" {
			oidcConfig.BackChannelLogoutUri = &val
		}
	}

	name := d.Get(NameVar).(string)

	_, err = client.UpdateApplication(ctx, &apppb.UpdateApplicationRequest{
		ApplicationId: d.Id(),
		ProjectId:     projectID,
		Name:          name,
		ApplicationType: &apppb.UpdateApplicationRequest_OidcConfiguration{
			OidcConfiguration: oidcConfig,
		},
	})
	if err != nil {
		return diag.Errorf("failed to update applicationOIDC: %v", err)
	}

	return nil
}

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

	respTypes := make([]apppb.OIDCResponseType, 0)
	for _, respType := range d.Get(responseTypesVar).([]interface{}) {
		respTypes = append(respTypes, apppb.OIDCResponseType(apppb.OIDCResponseType_value[respType.(string)]))
	}
	grantTypes := make([]apppb.OIDCGrantType, 0)
	for _, grantType := range d.Get(grantTypesVar).([]interface{}) {
		grantTypes = append(grantTypes, apppb.OIDCGrantType(apppb.OIDCGrantType_value[grantType.(string)]))
	}

	dur, err := time.ParseDuration(d.Get(clockSkewVar).(string))
	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := client.CreateApplication(ctx, &apppb.CreateApplicationRequest{
		ProjectId: d.Get(ProjectIDVar).(string),
		Name:      d.Get(NameVar).(string),
		ApplicationType: &apppb.CreateApplicationRequest_OidcConfiguration{
			OidcConfiguration: &apppb.CreateOIDCApplicationRequest{
				RedirectUris:             interfaceToStringSlice(d.Get(redirectURIsVar)),
				ResponseTypes:            respTypes,
				GrantTypes:               grantTypes,
				ApplicationType:          apppb.OIDCApplicationType(apppb.OIDCApplicationType_value[(d.Get(appTypeVar).(string))]),
				AuthMethodType:           apppb.OIDCAuthMethodType(apppb.OIDCAuthMethodType_value[(d.Get(authMethodTypeVar).(string))]),
				PostLogoutRedirectUris:   interfaceToStringSlice(d.Get(postLogoutRedirectURIsVar)),
				Version:                  apppb.OIDCVersion(apppb.OIDCVersion_value[d.Get(versionVar).(string)]),
				DevelopmentMode:          d.Get(devModeVar).(bool),
				AccessTokenType:          apppb.OIDCTokenType(apppb.OIDCTokenType_value[(d.Get(accessTokenTypeVar).(string))]),
				AccessTokenRoleAssertion: d.Get(accessTokenRoleAssertionVar).(bool),
				IdTokenRoleAssertion:     d.Get(idTokenRoleAssertionVar).(bool),
				IdTokenUserinfoAssertion: d.Get(idTokenUserinfoAssertionVar).(bool),
				ClockSkew:                durationpb.New(dur),
				AdditionalOrigins:        interfaceToStringSlice(d.Get(additionalOriginsVar)),
				SkipNativeAppSuccessPage: d.Get(skipNativeAppSuccessPageVar).(bool),
				BackChannelLogoutUri:     d.Get(BackChannelLogoutURIVar).(string),
				LoginVersion:             getLoginVersion(d),
			},
		},
	})
	if err != nil {
		return diag.Errorf("failed to create applicationOIDC: %v", err)
	}

	oidcConfig := resp.GetOidcConfiguration()
	set := map[string]interface{}{
		ClientIDVar:     oidcConfig.GetClientId(),
		ClientSecretVar: oidcConfig.GetClientSecret(),
	}
	for k, v := range set {
		if err := d.Set(k, v); err != nil {
			return diag.Errorf("failed to set %s of applicationOIDC: %v", k, err)
		}
	}

	d.SetId(resp.GetApplicationId())
	return read(ctx, d, m)
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

	resp, err := client.GetApplication(ctx, &apppb.GetApplicationRequest{ApplicationId: helper.GetID(d, AppIDVar)})
	if err != nil && helper.IgnoreIfNotFoundError(err) == nil {
		d.SetId("")
		return nil
	}
	if err != nil {
		return diag.Errorf("failed to get application oidc: %v", err)
	}

	app := resp.GetApplication()
	oidc := app.GetOidcConfiguration()
	grantTypes := make([]string, 0)
	for _, grantType := range oidc.GetGrantTypes() {
		grantTypes = append(grantTypes, grantType.String())
	}
	responseTypes := make([]string, 0)
	for _, responseType := range oidc.GetResponseTypes() {
		responseTypes = append(responseTypes, responseType.String())
	}
	clockSkew := oidc.GetClockSkew().AsDuration().String()
	if clockSkew == "" {
		clockSkew = "0s"
	}

	complianceProblems := make([]interface{}, 0)
	for _, p := range oidc.GetComplianceProblems() {
		complianceProblems = append(complianceProblems, map[string]interface{}{
			ComplianceProblemKeyVar:     p.GetKey(),
			ComplianceProblemMessageVar: p.GetLocalizedMessage(),
		})
	}

	loginVersion := []interface{}{}
	if oidc.GetLoginVersion() != nil {
		switch oidc.GetLoginVersion().GetVersion().(type) {
		case *apppb.LoginVersion_LoginV1:
			loginVersion = append(loginVersion, map[string]interface{}{
				LoginV1Var: true,
			})
		case *apppb.LoginVersion_LoginV2:
			v2 := oidc.GetLoginVersion().GetLoginV2()
			v2Map := map[string]interface{}{}
			if baseUri := v2.GetBaseUri(); baseUri != "" {
				v2Map[BaseURIVar] = baseUri
			}
			loginVersion = append(loginVersion, map[string]interface{}{
				LoginV2Var: []interface{}{v2Map},
			})
		}
	}

	set := map[string]interface{}{
		NameVar:                     app.GetName(),
		redirectURIsVar:             oidc.GetRedirectUris(),
		responseTypesVar:            responseTypes,
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
		ClientIDVar:                 oidc.GetClientId(),
		skipNativeAppSuccessPageVar: oidc.GetSkipNativeAppSuccessPage(),
		NoneCompliantVar:            oidc.GetNonCompliant(),
		ComplianceProblemsVar:       complianceProblems,
		BackChannelLogoutURIVar:     oidc.GetBackChannelLogoutUri(),
	}

	if len(loginVersion) > 0 {
		set[LoginVersionVar] = loginVersion
	}

	for k, v := range set {
		if err := d.Set(k, v); err != nil {
			return diag.Errorf("failed to set %s of applicationOIDC: %v", k, err)
		}
	}
	d.SetId(app.GetApplicationId())
	return nil
}

func interfaceToStringSlice(in interface{}) []string {
	slice := in.([]interface{})
	ret := make([]string, 0)
	for _, item := range slice {
		ret = append(ret, item.(string))
	}
	return ret
}

func getLoginVersion(d *schema.ResourceData) *apppb.LoginVersion {
	v, ok := d.GetOk(LoginVersionVar)
	if !ok {
		return nil
	}

	list := v.([]interface{})
	if len(list) == 0 {
		return nil
	}

	if list[0] == nil {
		return nil
	}

	item := list[0].(map[string]interface{})

	if loginV1, ok := item[LoginV1Var]; ok && loginV1.(bool) {
		return &apppb.LoginVersion{
			Version: &apppb.LoginVersion_LoginV1{
				LoginV1: &apppb.LoginV1{},
			},
		}
	}

	if v2, ok := item[LoginV2Var]; ok && v2 != nil {
		v2List := v2.([]interface{})
		if len(v2List) > 0 {
			if v2List[0] == nil {
				return nil
			}
			v2Item := v2List[0].(map[string]interface{})
			var uri *string
			if baseURI, ok := v2Item[BaseURIVar]; ok && baseURI.(string) != "" {
				uriStr := baseURI.(string)
				uri = &uriStr
			}
			return &apppb.LoginVersion{
				Version: &apppb.LoginVersion_LoginV2{
					LoginV2: &apppb.LoginV2{
						BaseUri: uri,
					},
				},
			}
		}
	}

	return nil
}

func list(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started list")
	name := d.Get(NameVar).(string)
	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}
	client, err := helper.GetAppV2Client(ctx, clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}
	projectID := d.Get(ProjectIDVar).(string)

	filters := make([]*apppb.ApplicationSearchFilter, 0)

	if projectID != "" {
		filters = append(filters, &apppb.ApplicationSearchFilter{
			Filter: &apppb.ApplicationSearchFilter_ProjectIdFilter{
				ProjectIdFilter: &apppb.ProjectIDFilter{
					ProjectId: projectID,
				},
			},
		})
	}

	if name != "" {
		filters = append(filters, &apppb.ApplicationSearchFilter{
			Filter: &apppb.ApplicationSearchFilter_NameFilter{
				NameFilter: &apppb.ApplicationNameFilter{
					Name: name,
				},
			},
		})
	}

	resp, err := client.ListApplications(ctx, &apppb.ListApplicationsRequest{
		Filters: filters,
	})
	if err != nil {
		return diag.Errorf("error while getting app by name %s: %v", name, err)
	}
	ids := make([]string, 0)
	for _, res := range resp.GetApplications() {
		if res.GetOidcConfiguration() == nil {
			continue
		}
		ids = append(ids, res.GetApplicationId())
	}
	d.SetId("-")
	return diag.FromErr(d.Set(appIDsVar, ids))
}
