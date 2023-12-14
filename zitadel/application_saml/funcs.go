package application_saml

import (
	"context"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/helper"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/management"
)

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

	request := &management.AddSAMLAppRequest{
		ProjectId: d.Get(ProjectIDVar).(string),
		Name:      d.Get(NameVar).(string),
	}

	if _, ok := d.GetOk(MetadataUrlVar); ok {
		request.Metadata = &management.AddSAMLAppRequest_MetadataUrl{
			MetadataUrl: d.Get(MetadataUrlVar).(string),
		}
	} else {
		request.Metadata = &management.AddSAMLAppRequest_MetadataXml{
			MetadataXml: d.Get(MetadataUrlVar).([]byte),
		}
	}

	response, err := client.AddSAMLApp(helper.CtxWithOrgID(ctx, d), request)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(response.GetAppId())
	return nil
}

func update(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started update")
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
		MetadataUrlVar,
		MetadataXmlVar,
	) {
		request := &management.UpdateSAMLAppConfigRequest{
			AppId:     d.Id(),
			ProjectId: projectID,
		}

		if _, ok := d.GetOk(MetadataUrlVar); ok {
			request.Metadata = &management.UpdateSAMLAppConfigRequest_MetadataUrl{
				MetadataUrl: d.Get(MetadataUrlVar).(string),
			}
		} else {
			request.Metadata = &management.UpdateSAMLAppConfigRequest_MetadataXml{
				MetadataXml: d.Get(MetadataUrlVar).([]byte),
			}
		}

		_, err := client.UpdateSAMLAppConfig(helper.CtxWithOrgID(ctx, d), request)

		if err != nil {
			return diag.Errorf("failed to update applicationSAML: %v", err)
		}
	}

	return nil
}

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

	resp, err := client.GetAppByID(helper.CtxWithOrgID(ctx, d), &management.GetAppByIDRequest{ProjectId: d.Get(ProjectIDVar).(string), AppId: helper.GetID(d, appIDVar)})
	if err != nil && helper.IgnoreIfNotFoundError(err) == nil {
		d.SetId("")
		return nil
	}
	if err != nil {
		return diag.Errorf("failed to get application oidc")
	}

	samlApp := resp.GetApp()
	saml := samlApp.GetSamlConfig()

	set := map[string]interface{}{
		helper.OrgIDVar: samlApp.GetDetails().GetResourceOwner(),
		NameVar:         samlApp.GetName(),
		MetadataXmlVar:  saml.GetMetadataXml(),
		MetadataUrlVar:  saml.GetMetadataUrl(),
	}

	for k, v := range set {
		if err := d.Set(k, v); err != nil {
			return diag.Errorf("failed to set %s of applicationSAML: %v", k, err)
		}
	}
	d.SetId(samlApp.GetId())

	return nil
}
