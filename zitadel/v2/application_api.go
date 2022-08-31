package v2

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zitadel/zitadel-go/v2/pkg/client/management"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/app"
	management2 "github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/management"
)

func GetApplicationAPI() *schema.Resource {
	return &schema.Resource{
		Description: "Resource representing an API application belonging to a project, with all configuration possibilities.",
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
			applicationAuthMethodTypeVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Auth method type",
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
		DeleteContext: deleteApplicationAPI,
		CreateContext: createApplicationAPI,
		UpdateContext: updateApplicationAPI,
		ReadContext:   readApplicationAPI,
		Importer:      &schema.ResourceImporter{StateContext: schema.ImportStatePassthroughContext},
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
		return diag.Errorf("failed to delete applicationAPI: %v", err)
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

	projectID := d.Get(applicationProjectIDVar).(string)
	appID := d.Id()
	apiApp, err := getApp(ctx, client, projectID, appID)

	appName := d.Get(applicationNameVar).(string)
	if apiApp.GetName() != appName {
		_, err = client.UpdateApp(ctx, &management2.UpdateAppRequest{
			ProjectId: projectID,
			AppId:     d.Id(),
			Name:      appName,
		})
		if err != nil {
			return diag.Errorf("failed to update application: %v", err)
		}
	}

	apiConfig := apiApp.GetApiConfig()
	authMethod := d.Get(applicationAuthMethodTypeVar).(string)
	if apiConfig.GetAuthMethodType().String() != authMethod {
		_, err = client.UpdateAPIAppConfig(ctx, &management2.UpdateAPIAppConfigRequest{
			ProjectId:      d.Get(applicationProjectIDVar).(string),
			AppId:          d.Id(),
			AuthMethodType: app.APIAuthMethodType(app.APIAuthMethodType_value[authMethod]),
		})
		if err != nil {
			return diag.Errorf("failed to update applicationAPI: %v", err)
		}
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

	set := map[string]interface{}{
		applicationClientID:     resp.GetClientId(),
		applicationClientSecret: resp.GetClientSecret(),
	}
	for k, v := range set {
		if err := d.Set(k, v); err != nil {
			return diag.Errorf("failed to set %s of applicationAPI: %v", k, err)
		}
	}
	if err != nil {
		return diag.Errorf("failed to create applicationAPI: %v", err)
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
		d.SetId("")
		return nil
		//return diag.FromErr(err)
	}

	app, err := getApp(ctx, client, d.Get(applicationProjectIDVar).(string), d.Id())
	if err != nil {
		return diag.Errorf("failed to read project: %v", err)
	}

	api := app.GetApiConfig()
	set := map[string]interface{}{
		applicationOrgIdVar:          app.GetDetails().GetResourceOwner(),
		applicationNameVar:           app.GetName(),
		applicationAuthMethodTypeVar: api.GetAuthMethodType().String(),
	}
	for k, v := range set {
		if err := d.Set(k, v); err != nil {
			return diag.Errorf("failed to set %s of applicationAPI: %v", k, err)
		}
	}
	d.SetId(app.GetId())
	return nil
}

func getApp(ctx context.Context, client *management.Client, projectID string, appID string) (*app.App, error) {
	resp, err := client.GetAppByID(ctx, &management2.GetAppByIDRequest{ProjectId: projectID, AppId: appID})
	if err != nil {
		return nil, fmt.Errorf("failed to read project: %v", err)
	}

	return resp.GetApp(), err
}
