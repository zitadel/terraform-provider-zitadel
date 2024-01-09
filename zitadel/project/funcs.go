package project

import (
	"context"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/management"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/object"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/project"

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

	_, err = client.RemoveProject(helper.CtxWithOrgID(ctx, d), &management.RemoveProjectRequest{
		Id: d.Id(),
	})
	if err != nil {
		return diag.Errorf("failed to delete project: %v", err)
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

	_, err = client.UpdateProject(helper.CtxWithOrgID(ctx, d), &management.UpdateProjectRequest{
		Id:                     d.Id(),
		Name:                   d.Get(NameVar).(string),
		ProjectRoleCheck:       d.Get(roleCheckVar).(bool),
		ProjectRoleAssertion:   d.Get(roleAssertionVar).(bool),
		HasProjectCheck:        d.Get(hasProjectCheckVar).(bool),
		PrivateLabelingSetting: project.PrivateLabelingSetting(project.PrivateLabelingSetting_value[d.Get(privateLabelingSettingVar).(string)]),
	})
	if err != nil {
		return diag.Errorf("failed to update project: %v", err)
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

	plSetting := d.Get(privateLabelingSettingVar).(string)
	resp, err := client.AddProject(helper.CtxWithOrgID(ctx, d), &management.AddProjectRequest{
		Name:                   d.Get(NameVar).(string),
		ProjectRoleAssertion:   d.Get(roleAssertionVar).(bool),
		ProjectRoleCheck:       d.Get(roleCheckVar).(bool),
		HasProjectCheck:        d.Get(hasProjectCheckVar).(bool),
		PrivateLabelingSetting: project.PrivateLabelingSetting(project.PrivateLabelingSetting_value[plSetting]),
	})
	if err != nil {
		return diag.Errorf("failed to create project: %v", err)
	}
	d.SetId(resp.GetId())
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

	resp, err := client.GetProjectByID(helper.CtxWithOrgID(ctx, d), &management.GetProjectByIDRequest{Id: helper.GetID(d, ProjectIDVar)})
	if err != nil && helper.IgnoreIfNotFoundError(err) == nil {
		d.SetId("")
		return nil
	}
	if err != nil {
		return diag.Errorf("failed to get project")
	}

	project := resp.GetProject()
	set := map[string]interface{}{
		helper.OrgIDVar:           project.GetDetails().GetResourceOwner(),
		stateVar:                  project.GetState().String(),
		NameVar:                   project.GetName(),
		roleAssertionVar:          project.GetProjectRoleAssertion(),
		roleCheckVar:              project.GetProjectRoleCheck(),
		hasProjectCheckVar:        project.GetHasProjectCheck(),
		privateLabelingSettingVar: project.PrivateLabelingSetting.String(),
	}
	for k, v := range set {
		if err := d.Set(k, v); err != nil {
			return diag.Errorf("failed to set %s of project: %v", k, err)
		}
	}
	d.SetId(project.GetId())

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
	client, err := helper.GetManagementClient(clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}
	req := &management.ListProjectsRequest{}
	if name != "" {
		req.Queries = append(req.Queries,
			&project.ProjectQuery{
				Query: &project.ProjectQuery_NameQuery{
					NameQuery: &project.ProjectNameQuery{
						Name:   name,
						Method: object.TextQueryMethod(object.TextQueryMethod_value[nameMethod]),
					},
				},
			})
	}

	resp, err := client.ListProjects(helper.CtxWithOrgID(ctx, d), req)
	if err != nil {
		return diag.Errorf("error while getting project by name %s: %v", name, err)
	}
	ids := make([]string, len(resp.Result))
	for i, res := range resp.Result {
		ids[i] = res.Id
	}
	// If the ID is blank, the datasource is deleted and not usable.
	d.SetId("-")
	return diag.FromErr(d.Set(projectIDsVar, ids))
}
