package project

import (
	"context"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/management"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/project"

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

	_, err = client.RemoveProject(ctx, &management.RemoveProjectRequest{
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

	client, err := helper.GetManagementClient(clientinfo, d.Get(orgIDVar).(string))
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.UpdateProject(ctx, &management.UpdateProjectRequest{
		Id:                     d.Id(),
		Name:                   d.Get(nameVar).(string),
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

	client, err := helper.GetManagementClient(clientinfo, d.Get(orgIDVar).(string))
	if err != nil {
		return diag.FromErr(err)
	}

	plSetting := d.Get(privateLabelingSettingVar).(string)
	resp, err := client.AddProject(ctx, &management.AddProjectRequest{
		Name:                   d.Get(nameVar).(string),
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

	client, err := helper.GetManagementClient(clientinfo, d.Get(orgIDVar).(string))
	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := client.GetProjectByID(ctx, &management.GetProjectByIDRequest{Id: helper.GetID(d, projectIDVar)})
	if err != nil && helper.IgnoreIfNotFoundError(err) == nil {
		d.SetId("")
		return nil
	}
	if err != nil {
		return diag.Errorf("failed to get project")
	}

	project := resp.GetProject()
	set := map[string]interface{}{
		orgIDVar:                  project.GetDetails().GetResourceOwner(),
		stateVar:                  project.GetState().String(),
		nameVar:                   project.GetName(),
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
