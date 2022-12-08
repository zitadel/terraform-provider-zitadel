package project_role

import (
	"context"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/management"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/object"
	project2 "github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/project"

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

	_, err = client.RemoveProjectRole(ctx, &management.RemoveProjectRoleRequest{
		ProjectId: d.Get(projectIDVar).(string),
		RoleKey:   d.Get(keyVar).(string),
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

	client, err := helper.GetManagementClient(clientinfo, d.Get(orgIDVar).(string))
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.UpdateProjectRole(ctx, &management.UpdateProjectRoleRequest{
		ProjectId:   d.Get(projectIDVar).(string),
		RoleKey:     d.Get(keyVar).(string),
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

	orgID := d.Get(orgIDVar).(string)
	client, err := helper.GetManagementClient(clientinfo, orgID)
	if err != nil {
		return diag.FromErr(err)
	}

	projectID := d.Get(projectIDVar).(string)
	roleKey := d.Get(keyVar).(string)
	_, err = client.AddProjectRole(ctx, &management.AddProjectRoleRequest{
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

	orgID := d.Get(orgIDVar).(string)
	client, err := helper.GetManagementClient(clientinfo, orgID)
	if err != nil {
		return diag.FromErr(err)
	}

	projectID := d.Get(projectIDVar).(string)
	resp, err := client.ListProjectRoles(ctx, &management.ListProjectRolesRequest{
		ProjectId: projectID,
		Queries: []*project2.RoleQuery{
			{Query: &project2.RoleQuery_KeyQuery{
				KeyQuery: &project2.RoleKeyQuery{
					Key:    d.Get(keyVar).(string),
					Method: object.TextQueryMethod_TEXT_QUERY_METHOD_EQUALS,
				},
			}},
		},
	})
	if err != nil || resp.Result == nil || len(resp.Result) == 0 {
		d.SetId("")
		return nil
	}

	if len(resp.Result) == 1 {
		projectRole := resp.GetResult()[0]
		roleKey := projectRole.GetKey()
		set := map[string]interface{}{
			projectIDVar:   projectID,
			orgIDVar:       orgID,
			keyVar:         roleKey,
			displayNameVar: projectRole.GetDisplayName(),
			groupVar:       projectRole.GetGroup(),
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
