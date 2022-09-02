package application_oidc

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	management2 "github.com/zitadel/zitadel-go/v2/pkg/client/management"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/app"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/management"
	"google.golang.org/protobuf/types/known/durationpb"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/helper"
)

func delete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started delete")

	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := helper.GetManagementClient(clientinfo, d.Get(orgIDVar).(string))
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.RemoveApp(ctx, &management.RemoveAppRequest{
		ProjectId: d.Get(projectIDVar).(string),
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

	client, err := helper.GetManagementClient(clientinfo, d.Get(orgIDVar).(string))
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

	projectID := d.Get(projectIDVar).(string)
	appID := d.Id()
	oidcApp, err := getApp(ctx, client, projectID, appID)
	if err != nil {
		return diag.FromErr(err)
	}

	appName := d.Get(nameVar).(string)
	if oidcApp.GetName() != appName {
		_, err = client.UpdateApp(ctx, &management.UpdateAppRequest{
			ProjectId: projectID,
			AppId:     appID,
			Name:      appName,
		})
		if err != nil {
			return diag.Errorf("failed to update application: %v", err)
		}
	}

	oidcConfig := oidcApp.GetOidcConfig()
	redirecURIs := interfaceToStringSlice(d.Get(redirectURIsVar))
	appType := d.Get(appTypeVar).(string)
	authMethodType := d.Get(authMethodTypeVar).(string)
	postLogoutRedirectURIs := interfaceToStringSlice(d.Get(postLogoutRedirectURIsVar))
	devMode := d.Get(devModeVar).(bool)
	accessTokenType := d.Get(accessTokenTypeVar).(string)
	accessTokenRoleAssertion := d.Get(accessTokenRoleAssertionVar).(bool)
	idTokenRoleAssertion := d.Get(idTokenRoleAssertionVar).(bool)
	idTokenUserinfoAssertion := d.Get(idTokenUserinfoAssertionVar).(bool)
	clockSkew := durationpb.New(dur)
	additionalOrigins := interfaceToStringSlice(d.Get(additionalOriginsVar))
	if !reflect.DeepEqual(redirecURIs, oidcConfig.GetRedirectUris()) ||
		!reflect.DeepEqual(respTypes, oidcConfig.GetResponseTypes()) ||
		!reflect.DeepEqual(grantTypes, oidcConfig.GetGrantTypes()) ||
		appType != oidcConfig.AppType.String() ||
		authMethodType != oidcConfig.AuthMethodType.String() ||
		!reflect.DeepEqual(postLogoutRedirectURIs, oidcConfig.GetPostLogoutRedirectUris()) ||
		devMode != oidcConfig.DevMode ||
		accessTokenType != oidcConfig.AccessTokenType.String() ||
		accessTokenRoleAssertion != oidcConfig.AccessTokenRoleAssertion ||
		clockSkew.String() != oidcConfig.ClockSkew.String() ||
		!reflect.DeepEqual(additionalOrigins, oidcConfig.GetAdditionalOrigins()) {
		_, err = client.UpdateOIDCAppConfig(ctx, &management.UpdateOIDCAppConfigRequest{
			ProjectId:                projectID,
			AppId:                    appID,
			RedirectUris:             redirecURIs,
			ResponseTypes:            respTypes,
			GrantTypes:               grantTypes,
			AppType:                  app.OIDCAppType(app.OIDCAppType_value[appType]),
			AuthMethodType:           app.OIDCAuthMethodType(app.OIDCAuthMethodType_value[authMethodType]),
			PostLogoutRedirectUris:   postLogoutRedirectURIs,
			DevMode:                  devMode,
			AccessTokenType:          app.OIDCTokenType(app.OIDCTokenType_value[accessTokenType]),
			AccessTokenRoleAssertion: accessTokenRoleAssertion,
			IdTokenRoleAssertion:     idTokenRoleAssertion,
			IdTokenUserinfoAssertion: idTokenUserinfoAssertion,
			ClockSkew:                clockSkew,
			AdditionalOrigins:        additionalOrigins,
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

	client, err := helper.GetManagementClient(clientinfo, d.Get(orgIDVar).(string))
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

	resp, err := client.AddOIDCApp(ctx, &management.AddOIDCAppRequest{
		ProjectId:                d.Get(projectIDVar).(string),
		Name:                     d.Get(nameVar).(string),
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
	})

	set := map[string]interface{}{
		clientID:     resp.GetClientId(),
		clientSecret: resp.GetClientSecret(),
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

	client, err := helper.GetManagementClient(clientinfo, d.Get(orgIDVar).(string))
	if err != nil {
		return diag.FromErr(err)
	}

	projectID := d.Get(projectIDVar).(string)
	oidcApp, err := getApp(ctx, client, projectID, helper.GetID(d, appIDVar))
	if err != nil {
		d.SetId("")
		return nil
		//return diag.Errorf("failed to read application: %v", err)
	}

	oidc := oidcApp.GetOidcConfig()
	grantTypes := make([]string, 0)
	for _, grantType := range oidc.GetGrantTypes() {
		grantTypes = append(grantTypes, grantType.String())
	}
	responseTypes := make([]string, 0)
	for _, responseType := range oidc.GetResponseTypes() {
		responseTypes = append(responseTypes, responseType.String())
	}
	clockSkew := oidc.GetClockSkew().String()
	if clockSkew == "" {
		clockSkew = "0s"
	}

	set := map[string]interface{}{
		orgIDVar:                    oidcApp.GetDetails().GetResourceOwner(),
		nameVar:                     oidcApp.GetName(),
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

func getApp(ctx context.Context, client *management2.Client, projectID string, appID string) (*app.App, error) {
	resp, err := client.GetAppByID(ctx, &management.GetAppByIDRequest{ProjectId: projectID, AppId: appID})
	if err != nil {
		return nil, fmt.Errorf("failed to read project: %v", err)
	}

	return resp.GetApp(), err
}
