package v2

import (
	"context"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/app"
	management2 "github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/management"
)

func GetApplicationAPI() *schema.Resource {
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
			applicationAuthMethodTypeVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Auth method type",
			},
		},
		DeleteContext: deleteApplicationAPI,
		CreateContext: createApplicationAPI,
		UpdateContext: updateApplicationAPI,
		ReadContext:   readApplicationAPI,
	}
}

func deleteApplicationAPI(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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

func updateApplicationAPI(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started update")

	clientinfo, ok := m.(*ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := getManagementClient(clientinfo, d.Get(applicationOrgIdVar).(string))
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
	_, err = client.UpdateAPIAppConfig(ctx, &management2.UpdateAPIAppConfigRequest{
		ProjectId:      d.Get(applicationProjectIDVar).(string),
		AppId:          d.Id(),
		AuthMethodType: app.APIAuthMethodType(app.APIAuthMethodType_value[(d.Get(applicationAuthMethodTypeVar).(string))]),
	})
	if err != nil {
		return diag.Errorf("failed to update applicationOIDC: %v", err)
	}
	return nil
}

func createApplicationAPI(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started create")

	clientinfo, ok := m.(*ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := getManagementClient(clientinfo, d.Get(applicationOrgIdVar).(string))
	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := client.AddAPIApp(ctx, &management2.AddAPIAppRequest{
		ProjectId:      d.Get(applicationProjectIDVar).(string),
		Name:           d.Get(applicationNameVar).(string),
		AuthMethodType: app.APIAuthMethodType(app.APIAuthMethodType_value[(d.Get(applicationAuthMethodTypeVar).(string))]),
	})

	if err != nil {
		return diag.Errorf("failed to create applicationOIDC: %v", err)
	}
	d.SetId(resp.GetAppId())
	return nil
}

func readApplicationAPI(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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
		return diag.Errorf("failed to read project: %v", err)
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
