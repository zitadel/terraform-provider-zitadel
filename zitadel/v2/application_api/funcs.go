package application_api

import (
	"context"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/app"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/management"

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
		return diag.Errorf("failed to delete applicationAPI: %v", err)
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

	projectID := d.Get(projectIDVar).(string)
	if d.HasChange(nameVar) {
		_, err = client.UpdateApp(ctx, &management.UpdateAppRequest{
			ProjectId: projectID,
			AppId:     d.Id(),
			Name:      d.Get(nameVar).(string),
		})
		if err != nil {
			return diag.Errorf("failed to update application: %v", err)
		}
	}

	if d.HasChanges(authMethodTypeVar) {
		_, err = client.UpdateAPIAppConfig(ctx, &management.UpdateAPIAppConfigRequest{
			ProjectId:      projectID,
			AppId:          d.Id(),
			AuthMethodType: app.APIAuthMethodType(app.APIAuthMethodType_value[d.Get(authMethodTypeVar).(string)]),
		})
		if err != nil {
			return diag.Errorf("failed to update applicationAPI: %v", err)
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

	resp, err := client.AddAPIApp(ctx, &management.AddAPIAppRequest{
		ProjectId:      d.Get(projectIDVar).(string),
		Name:           d.Get(nameVar).(string),
		AuthMethodType: app.APIAuthMethodType(app.APIAuthMethodType_value[(d.Get(authMethodTypeVar).(string))]),
	})

	set := map[string]interface{}{
		clientID:     resp.GetClientId(),
		clientSecret: resp.GetClientSecret(),
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

	resp, err := client.GetAppByID(ctx, &management.GetAppByIDRequest{ProjectId: d.Get(projectIDVar).(string), AppId: helper.GetID(d, appIDVar)})
	if err != nil && helper.IgnoreIfNotFoundError(err) == nil {
		d.SetId("")
		return nil
	}
	if err != nil {
		return diag.Errorf("failed to get application api")
	}

	app := resp.GetApp()
	api := app.GetApiConfig()
	set := map[string]interface{}{
		orgIDVar:          app.GetDetails().GetResourceOwner(),
		nameVar:           app.GetName(),
		authMethodTypeVar: api.GetAuthMethodType().String(),
	}
	for k, v := range set {
		if err := d.Set(k, v); err != nil {
			return diag.Errorf("failed to set %s of applicationAPI: %v", k, err)
		}
	}
	d.SetId(app.GetId())
	return nil
}
