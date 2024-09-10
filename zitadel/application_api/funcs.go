package application_api

import (
	"context"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/app"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/management"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/object"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
)

func delete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started delete")

	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := helper.GetManagementClient(ctx, clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.RemoveApp(helper.CtxWithOrgID(ctx, d), &management.RemoveAppRequest{
		ProjectId: d.Get(ProjectIDVar).(string),
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

	client, err := helper.GetManagementClient(ctx, clientinfo)
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

	if d.HasChanges(authMethodTypeVar) {
		_, err = client.UpdateAPIAppConfig(helper.CtxWithOrgID(ctx, d), &management.UpdateAPIAppConfigRequest{
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

	client, err := helper.GetManagementClient(ctx, clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := client.AddAPIApp(helper.CtxWithOrgID(ctx, d), &management.AddAPIAppRequest{
		ProjectId:      d.Get(ProjectIDVar).(string),
		Name:           d.Get(NameVar).(string),
		AuthMethodType: app.APIAuthMethodType(app.APIAuthMethodType_value[(d.Get(authMethodTypeVar).(string))]),
	})

	set := map[string]interface{}{
		ClientIDVar:     resp.GetClientId(),
		ClientSecretVar: resp.GetClientSecret(),
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

	client, err := helper.GetManagementClient(ctx, clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := client.GetAppByID(helper.CtxWithOrgID(ctx, d), &management.GetAppByIDRequest{ProjectId: d.Get(ProjectIDVar).(string), AppId: helper.GetID(d, AppIDVar)})
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
		helper.OrgIDVar:   app.GetDetails().GetResourceOwner(),
		NameVar:           app.GetName(),
		authMethodTypeVar: api.GetAuthMethodType().String(),
		ClientIDVar:       api.GetClientId(),
	}
	for k, v := range set {
		if err := d.Set(k, v); err != nil {
			return diag.Errorf("failed to set %s of applicationAPI: %v", k, err)
		}
	}
	d.SetId(app.GetId())
	return nil
}

func list(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started list")
	name := d.Get(NameVar).(string)
	nameMethod := d.Get(nameMethodVar).(string)
	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}
	client, err := helper.GetManagementClient(ctx, clientinfo)
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
		if res.GetApiConfig() == nil {
			continue
		}
		ids[i] = res.Id
	}
	// If the ID is blank, the datasource is deleted and not usable.
	d.SetId("-")
	return diag.FromErr(d.Set(appIDsVar, ids))
}
