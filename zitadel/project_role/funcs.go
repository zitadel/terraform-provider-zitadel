package project_role

import (
	"context"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/management"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/object"
	project2 "github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/project"

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

	_, err = client.RemoveProjectRole(helper.CtxWithOrgID(ctx, d), &management.RemoveProjectRoleRequest{
		ProjectId: d.Get(ProjectIDVar).(string),
		RoleKey:   d.Get(KeyVar).(string),
	})
	if err != nil {
		return diag.Errorf("failed to delete project role: %v", err)
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

	_, err = client.UpdateProjectRole(helper.CtxWithOrgID(ctx, d), &management.UpdateProjectRoleRequest{
		ProjectId:   d.Get(ProjectIDVar).(string),
		RoleKey:     d.Get(KeyVar).(string),
		DisplayName: d.Get(displayNameVar).(string),
		Group:       d.Get(groupVar).(string),
	})
	if err != nil {
		return diag.Errorf("failed to update project role: %v", err)
	}

	return nil
}

func create(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started create")

	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	orgID := d.Get(helper.OrgIDVar).(string)
	client, err := helper.GetManagementClient(ctx, clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}

	projectID := d.Get(ProjectIDVar).(string)
	roleKey := d.Get(KeyVar).(string)
	_, err = client.AddProjectRole(helper.CtxWithOrgID(ctx, d), &management.AddProjectRoleRequest{
		ProjectId:   projectID,
		RoleKey:     roleKey,
		DisplayName: d.Get(displayNameVar).(string),
		Group:       d.Get(groupVar).(string),
	})
	if err != nil {
		return diag.Errorf("failed to create project role: %v", err)
	}
	d.SetId(getProjectRoleID(orgID, projectID, roleKey))

	return nil
}

func read(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started read")

	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	orgID := d.Get(helper.OrgIDVar).(string)
	client, err := helper.GetManagementClient(ctx, clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}

	projectID := d.Get(ProjectIDVar).(string)
	resp, err := client.ListProjectRoles(helper.CtxWithOrgID(ctx, d), &management.ListProjectRolesRequest{
		ProjectId: projectID,
		Queries: []*project2.RoleQuery{
			{Query: &project2.RoleQuery_KeyQuery{
				KeyQuery: &project2.RoleKeyQuery{
					Key:    d.Get(KeyVar).(string),
					Method: object.TextQueryMethod_TEXT_QUERY_METHOD_EQUALS,
				},
			}},
		},
	})
	if err != nil {
		return diag.Errorf("failed to list project roles")
	}

	if len(resp.Result) == 1 {
		projectRole := resp.GetResult()[0]
		roleKey := projectRole.GetKey()
		set := map[string]interface{}{
			ProjectIDVar:    projectID,
			helper.OrgIDVar: orgID,
			KeyVar:          roleKey,
			displayNameVar:  projectRole.GetDisplayName(),
			groupVar:        projectRole.GetGroup(),
		}
		for k, v := range set {
			if err := d.Set(k, v); err != nil {
				return diag.Errorf("failed to set %s of project: %v", k, err)
			}
		}
		d.SetId(getProjectRoleID(orgID, projectID, roleKey))
		return nil
	}

	d.SetId("")
	return nil
}

func getProjectRoleID(orgID string, projectID string, roleKey string) string {
	return orgID + "_" + projectID + "_" + roleKey
}
