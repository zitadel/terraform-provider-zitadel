package application_saml

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
		return diag.Errorf("failed to delete applicationSAML: %v", err)
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

	if d.HasChanges(MetadataXMLVar, MetadataURLVar, LoginVersionVar) {
		req := &management.UpdateSAMLAppConfigRequest{
			ProjectId:    projectID,
			AppId:        d.Id(),
			LoginVersion: getLoginVersion(d),
		}
		if v, ok := d.GetOk(MetadataURLVar); ok && v.(string) != "" {
			req.Metadata = &management.UpdateSAMLAppConfigRequest_MetadataUrl{
				MetadataUrl: v.(string),
			}
		} else {
			req.Metadata = &management.UpdateSAMLAppConfigRequest_MetadataXml{
				MetadataXml: []byte(d.Get(MetadataXMLVar).(string)),
			}
		}
		_, err = client.UpdateSAMLAppConfig(helper.CtxWithOrgID(ctx, d), req)
		if err != nil {
			return diag.Errorf("failed to update applicationSAML: %v", err)
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

	req := &management.AddSAMLAppRequest{
		ProjectId:    d.Get(ProjectIDVar).(string),
		Name:         d.Get(NameVar).(string),
		LoginVersion: getLoginVersion(d),
	}
	if v, ok := d.GetOk(MetadataURLVar); ok && v.(string) != "" {
		req.Metadata = &management.AddSAMLAppRequest_MetadataUrl{
			MetadataUrl: v.(string),
		}
	} else {
		req.Metadata = &management.AddSAMLAppRequest_MetadataXml{
			MetadataXml: []byte(d.Get(MetadataXMLVar).(string)),
		}
	}

	resp, err := client.AddSAMLApp(helper.CtxWithOrgID(ctx, d), req)
	if err != nil {
		return diag.Errorf("failed to create applicationSAML: %v", err)
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
		return diag.Errorf("failed to get application saml: %v", err)
	}

	samlApp := resp.GetApp()
	samlConfig := samlApp.GetSamlConfig()
	set := map[string]interface{}{
		helper.OrgIDVar: samlApp.GetDetails().GetResourceOwner(),
		NameVar:         samlApp.GetName(),
	}
	// Only set metadata_xml if the user did not configure metadata_url,
	// otherwise the resolved XML would cause a perpetual diff.
	if _, urlSet := d.GetOk(MetadataURLVar); !urlSet {
		set[MetadataXMLVar] = string(samlConfig.GetMetadataXml())
	}

	loginVersion := []interface{}{}
	if samlConfig.GetLoginVersion() != nil {
		switch samlConfig.GetLoginVersion().GetVersion().(type) {
		case *app.LoginVersion_LoginV1:
			loginVersion = append(loginVersion, map[string]interface{}{
				LoginV1Var: true,
			})
		case *app.LoginVersion_LoginV2:
			v2 := samlConfig.GetLoginVersion().GetLoginV2()
			v2Map := map[string]interface{}{}

			if baseUri := v2.GetBaseUri(); baseUri != "" {
				v2Map[BaseURIVar] = baseUri
			}

			loginVersion = append(loginVersion, map[string]interface{}{
				LoginV2Var: []interface{}{v2Map},
			})
		}
	}

	if len(loginVersion) > 0 {
		set[LoginVersionVar] = loginVersion
	}

	for k, v := range set {
		if err := d.Set(k, v); err != nil {
			return diag.Errorf("failed to set %s of applicationSAML: %v", k, err)
		}
	}
	d.SetId(samlApp.GetId())
	return nil
}

func getLoginVersion(d *schema.ResourceData) *app.LoginVersion {
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
		return &app.LoginVersion{
			Version: &app.LoginVersion_LoginV1{
				LoginV1: &app.LoginV1{},
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
			return &app.LoginVersion{
				Version: &app.LoginVersion_LoginV2{
					LoginV2: &app.LoginV2{
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
		if res.GetSamlConfig() == nil {
			continue
		}
		ids[i] = res.Id
	}
	// If the ID is blank, the datasource is deleted and not usable.
	d.SetId("-")
	return diag.FromErr(d.Set(appIDsVar, ids))
}
