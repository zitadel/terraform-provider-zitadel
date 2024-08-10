package application_oidc

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/app"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/management"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/object"
	"google.golang.org/protobuf/types/known/durationpb"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/helper"
)

func delete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started delete")

	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := helper.GetManagementClient(clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.RemoveApp(helper.CtxWithOrgID(ctx, d), &management.RemoveAppRequest{
		ProjectId: d.Get(ProjectIDVar).(string),
		AppId:     d.Id(),
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

	client, err := helper.GetManagementClient(clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}

	projectID := d.Get(ProjectIDVar).(string)

	if d.HasChange(NameVar) {
		_, err = client.UpdateApp(helper.CtxWithOrgID(ctx, d), &management.UpdateAppRequest{
			ProjectId: projectID,
			AppId:     d.Id(),
			Name:      d.Get(NameVar).(string),
		})
		if err != nil {
			return diag.Errorf("failed to update application: %v", err)
		}
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
	) {
		respTypes := make([]app.OIDCResponseType, 0)
		for _, respType := range d.Get(responseTypesVar).([]interface{}) {
			respTypes = append(respTypes, app.OIDCResponseType(app.OIDCResponseType_value[respType.(string)]))
		}
		grantTypes := make([]app.OIDCGrantType, 0)
		for _, grantType := range d.Get(grantTypesVar).([]interface{}) {
			grantTypes = append(grantTypes, app.OIDCGrantType(app.OIDCGrantType_value[grantType.(string)]))
		}
		dur, err := time.ParseDuration(d.Get(clockSkewVar).(string))
		if err != nil {
			return diag.FromErr(err)
		}

		_, err = client.UpdateOIDCAppConfig(helper.CtxWithOrgID(ctx, d), &management.UpdateOIDCAppConfigRequest{
			ProjectId:                projectID,
			AppId:                    d.Id(),
			RedirectUris:             interfaceToStringSlice(d.Get(redirectURIsVar)),
			ResponseTypes:            respTypes,
			GrantTypes:               grantTypes,
			AppType:                  app.OIDCAppType(app.OIDCAppType_value[d.Get(appTypeVar).(string)]),
			AuthMethodType:           app.OIDCAuthMethodType(app.OIDCAuthMethodType_value[d.Get(authMethodTypeVar).(string)]),
			PostLogoutRedirectUris:   interfaceToStringSlice(d.Get(postLogoutRedirectURIsVar)),
			DevMode:                  d.Get(devModeVar).(bool),
			AccessTokenType:          app.OIDCTokenType(app.OIDCTokenType_value[d.Get(accessTokenTypeVar).(string)]),
			AccessTokenRoleAssertion: d.Get(accessTokenRoleAssertionVar).(bool),
			IdTokenRoleAssertion:     d.Get(idTokenRoleAssertionVar).(bool),
			IdTokenUserinfoAssertion: d.Get(idTokenUserinfoAssertionVar).(bool),
			AdditionalOrigins:        interfaceToStringSlice(d.Get(additionalOriginsVar)),
			ClockSkew:                durationpb.New(dur),
		})
		if err != nil {
			return diag.Errorf("failed to update applicationOIDC: %v", err)
		}
	}
	return nil
}

func create(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started create")

	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := helper.GetManagementClient(clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}

	respTypes := make([]app.OIDCResponseType, 0)
	for _, respType := range d.Get(responseTypesVar).([]interface{}) {
		respTypes = append(respTypes, app.OIDCResponseType(app.OIDCResponseType_value[respType.(string)]))
	}
	grantTypes := make([]app.OIDCGrantType, 0)
	for _, grantType := range d.Get(grantTypesVar).([]interface{}) {
		grantTypes = append(grantTypes, app.OIDCGrantType(app.OIDCGrantType_value[grantType.(string)]))
	}

	dur, err := time.ParseDuration(d.Get(clockSkewVar).(string))
	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := client.AddOIDCApp(helper.CtxWithOrgID(ctx, d), &management.AddOIDCAppRequest{
		ProjectId:                d.Get(ProjectIDVar).(string),
		Name:                     d.Get(NameVar).(string),
		RedirectUris:             interfaceToStringSlice(d.Get(redirectURIsVar)),
		ResponseTypes:            respTypes,
		GrantTypes:               grantTypes,
		AppType:                  app.OIDCAppType(app.OIDCAppType_value[(d.Get(appTypeVar).(string))]),
		AuthMethodType:           app.OIDCAuthMethodType(app.OIDCAuthMethodType_value[(d.Get(authMethodTypeVar).(string))]),
		PostLogoutRedirectUris:   interfaceToStringSlice(d.Get(postLogoutRedirectURIsVar)),
		DevMode:                  d.Get(devModeVar).(bool),
		AccessTokenType:          app.OIDCTokenType(app.OIDCTokenType_value[(d.Get(accessTokenTypeVar).(string))]),
		AccessTokenRoleAssertion: d.Get(accessTokenRoleAssertionVar).(bool),
		IdTokenRoleAssertion:     d.Get(idTokenRoleAssertionVar).(bool),
		IdTokenUserinfoAssertion: d.Get(idTokenUserinfoAssertionVar).(bool),
		ClockSkew:                durationpb.New(dur),
		AdditionalOrigins:        interfaceToStringSlice(d.Get(additionalOriginsVar)),
		Version:                  app.OIDCVersion(app.OIDCVersion_value[d.Get(versionVar).(string)]),
	})

	set := map[string]interface{}{
		ClientIDVar:     resp.GetClientId(),
		ClientSecretVar: resp.GetClientSecret(),
	}
	for k, v := range set {
		if err := d.Set(k, v); err != nil {
			return diag.Errorf("failed to set %s of applicationOIDC: %v", k, err)
		}
	}

	if err != nil {
		return diag.Errorf("failed to create applicationOIDC: %v", err)
	}
	d.SetId(resp.GetAppId())
	return nil
}

func read(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started read")

	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := helper.GetManagementClient(clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := client.GetAppByID(helper.CtxWithOrgID(ctx, d), &management.GetAppByIDRequest{ProjectId: d.Get(ProjectIDVar).(string), AppId: helper.GetID(d, AppIDVar)})
	if err != nil && helper.IgnoreIfNotFoundError(err) == nil {
		d.SetId("")
		return nil
	}
	if err != nil {
		return diag.Errorf("failed to get application oidc")
	}

	oidcApp := resp.GetApp()
	oidc := oidcApp.GetOidcConfig()
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

	set := map[string]interface{}{
		helper.OrgIDVar:             oidcApp.GetDetails().GetResourceOwner(),
		NameVar:                     oidcApp.GetName(),
		redirectURIsVar:             oidc.GetRedirectUris(),
		responseTypesVar:            responseTypes,
		grantTypesVar:               grantTypes,
		appTypeVar:                  oidc.GetAppType().String(),
		authMethodTypeVar:           oidc.GetAuthMethodType().String(),
		postLogoutRedirectURIsVar:   oidc.GetPostLogoutRedirectUris(),
		versionVar:                  oidc.GetVersion().String(),
		devModeVar:                  oidc.GetDevMode(),
		accessTokenTypeVar:          oidc.GetAccessTokenType().String(),
		accessTokenRoleAssertionVar: oidc.GetAccessTokenRoleAssertion(),
		idTokenRoleAssertionVar:     oidc.GetIdTokenRoleAssertion(),
		idTokenUserinfoAssertionVar: oidc.GetIdTokenUserinfoAssertion(),
		clockSkewVar:                clockSkew,
		additionalOriginsVar:        oidc.GetAdditionalOrigins(),
		ClientIDVar:                 oidc.GetClientId(),
	}
	for k, v := range set {
		if err := d.Set(k, v); err != nil {
			return diag.Errorf("failed to set %s of applicationOIDC: %v", k, err)
		}
	}
	d.SetId(oidcApp.GetId())
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

func list(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started list")
	name := d.Get(NameVar).(string)
	nameMethod := d.Get(nameMethodVar).(string)
	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}
	client, err := helper.GetManagementClient(clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}
	req := &management.ListAppsRequest{
		ProjectId: d.Get(ProjectIDVar).(string),
	}
	if name != "" {
		req.Queries = append(req.Queries,
			&app.AppQuery{
				Query: &app.AppQuery_NameQuery{
					NameQuery: &app.AppNameQuery{
						Name:   name,
						Method: object.TextQueryMethod(object.TextQueryMethod_value[nameMethod]),
					},
				},
			})
	}
	resp, err := client.ListApps(helper.CtxWithOrgID(ctx, d), req)
	if err != nil {
		return diag.Errorf("error while getting app by name %s: %v", name, err)
	}
	ids := make([]string, len(resp.Result))
	for i, res := range resp.Result {
		if res.GetOidcConfig() == nil {
			continue
		}
		ids[i] = res.Id
	}
	// If the ID is blank, the datasource is deleted and not usable.
	d.SetId("-")
	return diag.FromErr(d.Set(appIDsVar, ids))
}
