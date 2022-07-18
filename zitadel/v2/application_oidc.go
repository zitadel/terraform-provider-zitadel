package v2

import (
	"context"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/app"
	management2 "github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/management"
	"google.golang.org/protobuf/types/known/durationpb"
	"time"
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
)

func GetApplicationOIDC() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			applicationOrgIdVar: {
				Type:        schema.TypeString,
				Computed:    true,
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
				Required:    true,
				Description: "Post logout redirect URIs",
			},
			applicationVersionVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Version",
			},
			applicationDevModeVar: {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Dev mode",
			},
			applicationAccessTokenTypeVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Access token type",
			},
			applicationAccessTokenRoleAssertionVar: {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Access token role assertion",
			},
			applicationIdTokenRoleAssertionVar: {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "ID token role assertion",
			},
			applicationIdTokenUserinfoAssertionVar: {
				Type:        schema.TypeBool,
				Required:    true,
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
				Required:    true,
				Description: "Additional origins",
			},
		},
		DeleteContext: deleteApplicationOIDC,
		CreateContext: createApplicationOIDC,
		UpdateContext: updateApplicationOIDC,
		ReadContext:   readApplicationOIDC,
	}
}

func deleteApplicationOIDC(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started delete")

	clientinfo, ok := m.(*ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := getManagementClient(clientinfo, d.Get(projectResourceOwner).(string))
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
	for _, respType := range d.Get(applicationResponseTypesVar).([]string) {
		respTypes = append(respTypes, app.OIDCResponseType(app.OIDCResponseType_value[respType]))
	}
	grantTypes := make([]app.OIDCGrantType, 0)
	for _, grantType := range d.Get(applicationGrantTypesVar).([]string) {
		grantTypes = append(grantTypes, app.OIDCGrantType(app.OIDCGrantType_value[grantType]))
	}

	dur, err := time.ParseDuration(d.Get(applicationClockSkewVar).(string))
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.UpdateApp(ctx, &management2.UpdateAppRequest{
		ProjectId: d.Get(applicationProjectIDVar).(string),
		AppId:     d.Id(),
		Name:      d.Get(applicationNameVar).(string),
	})
	if err != nil {
		return diag.Errorf("failed to update application: %v", err)
	}
	_, err = client.UpdateOIDCAppConfig(ctx, &management2.UpdateOIDCAppConfigRequest{
		ProjectId:                d.Get(applicationProjectIDVar).(string),
		AppId:                    d.Id(),
		RedirectUris:             d.Get(applicationRedirectURIsVar).([]string),
		ResponseTypes:            respTypes,
		GrantTypes:               grantTypes,
		AppType:                  app.OIDCAppType(app.OIDCAppType_value[(d.Get(applicationAppTypeVar).(string))]),
		AuthMethodType:           app.OIDCAuthMethodType(app.OIDCAuthMethodType_value[(d.Get(applicationAuthMethodTypeVar).(string))]),
		PostLogoutRedirectUris:   d.Get(applicationPostLogoutRedirectURIsVar).([]string),
		DevMode:                  d.Get(applicationDevModeVar).(bool),
		AccessTokenType:          app.OIDCTokenType(app.OIDCTokenType_value[(d.Get(applicationAccessTokenTypeVar).(string))]),
		AccessTokenRoleAssertion: d.Get(applicationAccessTokenRoleAssertionVar).(bool),
		IdTokenRoleAssertion:     d.Get(applicationIdTokenRoleAssertionVar).(bool),
		IdTokenUserinfoAssertion: d.Get(applicationIdTokenUserinfoAssertionVar).(bool),
		ClockSkew:                durationpb.New(dur),
		AdditionalOrigins:        d.Get(applicationAdditionalOriginsVar).([]string),
	})
	if err != nil {
		return diag.Errorf("failed to update applicationOIDC: %v", err)
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
	for _, respType := range d.Get(applicationResponseTypesVar).([]string) {
		respTypes = append(respTypes, app.OIDCResponseType(app.OIDCResponseType_value[respType]))
	}
	grantTypes := make([]app.OIDCGrantType, 0)
	for _, grantType := range d.Get(applicationGrantTypesVar).([]string) {
		grantTypes = append(grantTypes, app.OIDCGrantType(app.OIDCGrantType_value[grantType]))
	}

	dur, err := time.ParseDuration(d.Get(applicationClockSkewVar).(string))
	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := client.AddOIDCApp(ctx, &management2.AddOIDCAppRequest{
		ProjectId:                d.Get(applicationProjectIDVar).(string),
		Name:                     d.Get(applicationNameVar).(string),
		RedirectUris:             d.Get(applicationRedirectURIsVar).([]string),
		ResponseTypes:            respTypes,
		GrantTypes:               grantTypes,
		AppType:                  app.OIDCAppType(app.OIDCAppType_value[(d.Get(applicationAppTypeVar).(string))]),
		AuthMethodType:           app.OIDCAuthMethodType(app.OIDCAuthMethodType_value[(d.Get(applicationAuthMethodTypeVar).(string))]),
		PostLogoutRedirectUris:   d.Get(applicationPostLogoutRedirectURIsVar).([]string),
		DevMode:                  d.Get(applicationDevModeVar).(bool),
		AccessTokenType:          app.OIDCTokenType(app.OIDCTokenType_value[(d.Get(applicationAccessTokenTypeVar).(string))]),
		AccessTokenRoleAssertion: d.Get(applicationAccessTokenRoleAssertionVar).(bool),
		IdTokenRoleAssertion:     d.Get(applicationIdTokenRoleAssertionVar).(bool),
		IdTokenUserinfoAssertion: d.Get(applicationIdTokenUserinfoAssertionVar).(bool),
		ClockSkew:                durationpb.New(dur),
		AdditionalOrigins:        d.Get(applicationAdditionalOriginsVar).([]string),
	})

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
		return diag.Errorf("failed to read application: %v", err)
	}

	app := resp.GetApp()
	oidc := app.GetOidcConfig()
	set := map[string]interface{}{
		applicationProjectIDVar:                app.GetDetails().GetResourceOwner(),
		applicationNameVar:                     app.GetName(),
		applicationRedirectURIsVar:             oidc.GetRedirectUris(),
		applicationResponseTypesVar:            oidc.GetResponseTypes(),
		applicationGrantTypesVar:               oidc.GetGrantTypes(),
		applicationAppTypeVar:                  oidc.GetAppType(),
		applicationAuthMethodTypeVar:           oidc.GetAuthMethodType(),
		applicationPostLogoutRedirectURIsVar:   oidc.GetPostLogoutRedirectUris(),
		applicationVersionVar:                  oidc.GetVersion(),
		applicationDevModeVar:                  oidc.GetDevMode(),
		applicationAccessTokenTypeVar:          oidc.GetAccessTokenType(),
		applicationAccessTokenRoleAssertionVar: oidc.GetAccessTokenRoleAssertion(),
		applicationIdTokenRoleAssertionVar:     oidc.GetIdTokenRoleAssertion(),
		applicationIdTokenUserinfoAssertionVar: oidc.GetIdTokenUserinfoAssertion(),
		applicationClockSkewVar:                oidc.GetClockSkew(),
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
