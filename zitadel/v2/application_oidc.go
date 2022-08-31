package v2

import (
	"context"
	"reflect"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/app"
	management2 "github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/management"
	"google.golang.org/protobuf/types/known/durationpb"
)

const (
	applicationOrgIdVar                    = "org_id"
	applicationProjectIDVar                = "project_id"
	applicationNameVar                     = "name"
	applicationRedirectURIsVar             = "redirect_uris"
	applicationResponseTypesVar            = "response_types"
	applicationGrantTypesVar               = "grant_types"
	applicationAppTypeVar                  = "app_type"
	applicationAuthMethodTypeVar           = "auth_method_type"
	applicationPostLogoutRedirectURIsVar   = "post_logout_redirect_uris"
	applicationVersionVar                  = "version"
	applicationDevModeVar                  = "dev_mode"
	applicationAccessTokenTypeVar          = "access_token_type"
	applicationAccessTokenRoleAssertionVar = "access_token_role_assertion"
	applicationIdTokenRoleAssertionVar     = "id_token_role_assertion"
	applicationIdTokenUserinfoAssertionVar = "id_token_userinfo_assertion"
	applicationClockSkewVar                = "clock_skew"
	applicationAdditionalOriginsVar        = "additional_origins"
	applicationClientID                    = "client_id"
	applicationClientSecret                = "client_secret"
)

func GetApplicationOIDC() *schema.Resource {
	return &schema.Resource{
		Description: "Resource representing an OIDC application belonging to a project, with all configuration possibilities.",
		Schema: map[string]*schema.Schema{
			applicationOrgIdVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "orgID of the application",
				ForceNew:    true,
			},
			applicationProjectIDVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the project",
				ForceNew:    true,
			},
			applicationNameVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the application",
			},
			applicationRedirectURIsVar: {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Required:    true,
				Description: "RedirectURIs",
			},
			applicationResponseTypesVar: {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Required:    true,
				Description: "Response type",
			},
			applicationGrantTypesVar: {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Required:    true,
				Description: "Grant types",
			},
			applicationAppTypeVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "App type",
			},
			applicationAuthMethodTypeVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Auth method type",
			},
			applicationPostLogoutRedirectURIsVar: {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				Description: "Post logout redirect URIs",
			},
			applicationVersionVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Version",
			},
			applicationDevModeVar: {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Dev mode",
			},
			applicationAccessTokenTypeVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Access token type",
			},
			applicationAccessTokenRoleAssertionVar: {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Access token role assertion",
			},
			applicationIdTokenRoleAssertionVar: {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "ID token role assertion",
			},
			applicationIdTokenUserinfoAssertionVar: {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Token userinfo assertion",
			},
			applicationClockSkewVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Clockskew",
			},
			applicationAdditionalOriginsVar: {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				Description: "Additional origins",
			},
			applicationClientID: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "generated ID for this config",
				Sensitive:   true,
			},
			applicationClientSecret: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "generated secret for this config",
				Sensitive:   true,
			},
		},
		DeleteContext: deleteApplicationOIDC,
		CreateContext: createApplicationOIDC,
		UpdateContext: updateApplicationOIDC,
		ReadContext:   readApplicationOIDC,
		Importer:      &schema.ResourceImporter{StateContext: schema.ImportStatePassthroughContext},
	}
}

func deleteApplicationOIDC(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started delete")

	clientinfo, ok := m.(*ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := getManagementClient(clientinfo, d.Get(applicationOrgIdVar).(string))
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.RemoveApp(ctx, &management2.RemoveAppRequest{
		ProjectId: d.Get(applicationProjectIDVar).(string),
		AppId:     d.Id(),
	})
	if err != nil {
		return diag.Errorf("failed to delete applicationOIDC: %v", err)
	}
	return nil
}

func updateApplicationOIDC(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started update")

	clientinfo, ok := m.(*ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := getManagementClient(clientinfo, d.Get(applicationOrgIdVar).(string))
	if err != nil {
		return diag.FromErr(err)
	}

	respTypes := make([]app.OIDCResponseType, 0)
	for _, respType := range d.Get(applicationResponseTypesVar).([]interface{}) {
		respTypes = append(respTypes, app.OIDCResponseType(app.OIDCResponseType_value[respType.(string)]))
	}
	grantTypes := make([]app.OIDCGrantType, 0)
	for _, grantType := range d.Get(applicationGrantTypesVar).([]interface{}) {
		grantTypes = append(grantTypes, app.OIDCGrantType(app.OIDCGrantType_value[grantType.(string)]))
	}

	dur, err := time.ParseDuration(d.Get(applicationClockSkewVar).(string))
	if err != nil {
		return diag.FromErr(err)
	}

	projectID := d.Get(applicationProjectIDVar).(string)
	appID := d.Id()
	apiApp, err := getApp(ctx, client, projectID, appID)

	appName := d.Get(applicationNameVar).(string)
	if apiApp.GetName() != appName {
		_, err = client.UpdateApp(ctx, &management2.UpdateAppRequest{
			ProjectId: projectID,
			AppId:     appID,
			Name:      appName,
		})
		if err != nil {
			return diag.Errorf("failed to update application: %v", err)
		}
	}

	oidcConfig := apiApp.GetOidcConfig()
	redirecURIs := interfaceToStringSlice(d.Get(applicationRedirectURIsVar))
	appType := d.Get(applicationAppTypeVar).(string)
	authMethodType := d.Get(applicationAuthMethodTypeVar).(string)
	postLogoutRedirectURIs := interfaceToStringSlice(d.Get(applicationPostLogoutRedirectURIsVar))
	devMode := d.Get(applicationDevModeVar).(bool)
	accessTokenType := d.Get(applicationAccessTokenTypeVar).(string)
	accessTokenRoleAssertion := d.Get(applicationAccessTokenRoleAssertionVar).(bool)
	idTokenRoleAssertion := d.Get(applicationIdTokenRoleAssertionVar).(bool)
	idTokenUserinfoAssertion := d.Get(applicationIdTokenUserinfoAssertionVar).(bool)
	clockSkew := durationpb.New(dur)
	additionalOrigins := interfaceToStringSlice(d.Get(applicationAdditionalOriginsVar))
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
		_, err = client.UpdateOIDCAppConfig(ctx, &management2.UpdateOIDCAppConfigRequest{
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

func createApplicationOIDC(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started create")

	clientinfo, ok := m.(*ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := getManagementClient(clientinfo, d.Get(applicationOrgIdVar).(string))
	if err != nil {
		return diag.FromErr(err)
	}

	respTypes := make([]app.OIDCResponseType, 0)
	for _, respType := range d.Get(applicationResponseTypesVar).([]interface{}) {
		respTypes = append(respTypes, app.OIDCResponseType(app.OIDCResponseType_value[respType.(string)]))
	}
	grantTypes := make([]app.OIDCGrantType, 0)
	for _, grantType := range d.Get(applicationGrantTypesVar).([]interface{}) {
		grantTypes = append(grantTypes, app.OIDCGrantType(app.OIDCGrantType_value[grantType.(string)]))
	}

	dur, err := time.ParseDuration(d.Get(applicationClockSkewVar).(string))
	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := client.AddOIDCApp(ctx, &management2.AddOIDCAppRequest{
		ProjectId:                d.Get(applicationProjectIDVar).(string),
		Name:                     d.Get(applicationNameVar).(string),
		RedirectUris:             interfaceToStringSlice(d.Get(applicationRedirectURIsVar)),
		ResponseTypes:            respTypes,
		GrantTypes:               grantTypes,
		AppType:                  app.OIDCAppType(app.OIDCAppType_value[(d.Get(applicationAppTypeVar).(string))]),
		AuthMethodType:           app.OIDCAuthMethodType(app.OIDCAuthMethodType_value[(d.Get(applicationAuthMethodTypeVar).(string))]),
		PostLogoutRedirectUris:   interfaceToStringSlice(d.Get(applicationPostLogoutRedirectURIsVar)),
		DevMode:                  d.Get(applicationDevModeVar).(bool),
		AccessTokenType:          app.OIDCTokenType(app.OIDCTokenType_value[(d.Get(applicationAccessTokenTypeVar).(string))]),
		AccessTokenRoleAssertion: d.Get(applicationAccessTokenRoleAssertionVar).(bool),
		IdTokenRoleAssertion:     d.Get(applicationIdTokenRoleAssertionVar).(bool),
		IdTokenUserinfoAssertion: d.Get(applicationIdTokenUserinfoAssertionVar).(bool),
		ClockSkew:                durationpb.New(dur),
		AdditionalOrigins:        interfaceToStringSlice(d.Get(applicationAdditionalOriginsVar)),
	})

	set := map[string]interface{}{
		applicationClientID:     resp.GetClientId(),
		applicationClientSecret: resp.GetClientSecret(),
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

func readApplicationOIDC(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started read")

	clientinfo, ok := m.(*ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := getManagementClient(clientinfo, d.Get(applicationOrgIdVar).(string))
	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := client.GetAppByID(ctx, &management2.GetAppByIDRequest{ProjectId: d.Get(applicationProjectIDVar).(string), AppId: d.Id()})
	if err != nil {
		d.SetId("")
		return nil
		//return diag.Errorf("failed to read application: %v", err)
	}

	app := resp.GetApp()
	oidc := app.GetOidcConfig()
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
		applicationOrgIdVar:                    app.GetDetails().GetResourceOwner(),
		applicationNameVar:                     app.GetName(),
		applicationRedirectURIsVar:             oidc.GetRedirectUris(),
		applicationResponseTypesVar:            responseTypes,
		applicationGrantTypesVar:               grantTypes,
		applicationAppTypeVar:                  oidc.GetAppType().String(),
		applicationAuthMethodTypeVar:           oidc.GetAuthMethodType().String(),
		applicationPostLogoutRedirectURIsVar:   oidc.GetPostLogoutRedirectUris(),
		applicationVersionVar:                  oidc.GetVersion().String(),
		applicationDevModeVar:                  oidc.GetDevMode(),
		applicationAccessTokenTypeVar:          oidc.GetAccessTokenType().String(),
		applicationAccessTokenRoleAssertionVar: oidc.GetAccessTokenRoleAssertion(),
		applicationIdTokenRoleAssertionVar:     oidc.GetIdTokenRoleAssertion(),
		applicationIdTokenUserinfoAssertionVar: oidc.GetIdTokenUserinfoAssertion(),
		applicationClockSkewVar:                clockSkew,
		applicationAdditionalOriginsVar:        oidc.GetAdditionalOrigins(),
	}
	for k, v := range set {
		if err := d.Set(k, v); err != nil {
			return diag.Errorf("failed to set %s of applicationOIDC: %v", k, err)
		}
	}
	d.SetId(app.GetId())
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
