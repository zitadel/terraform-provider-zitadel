package v2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	management2 "github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/management"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/object"
	project2 "github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/project"
)

const (
	projectRoleOrgID       = "org_id"
	projectRoleProjectID   = "project_id"
	projectRoleKey         = "role_key"
	projectRoleDisplayName = "display_name"
	projectRoleGroup       = "group"
)

func GetProjectRole() *schema.Resource {
	return &schema.Resource{
		Description: "Resource representing the project roles, which can be given as authorizations to users.",
		Schema: map[string]*schema.Schema{
			projectRoleProjectID: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the project",
				ForceNew:    true,
			},
			projectRoleOrgID: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the organization",
				ForceNew:    true,
			},
			projectRoleKey: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Key used for project role",
			},
			projectRoleDisplayName: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name used for project role",
			},
			projectRoleGroup: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Group used for project role",
			},
		},
		DeleteContext: deleteProjectRole,
		CreateContext: createProjectRole,
		UpdateContext: updateProjectRole,
		ReadContext:   readProjectRole,
	}
}

func deleteProjectRole(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started delete")

	clientinfo, ok := m.(*ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := getManagementClient(clientinfo, d.Get(projectRoleOrgID).(string))
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.RemoveProjectRole(ctx, &management2.RemoveProjectRoleRequest{
		ProjectId: d.Get(projectRoleProjectID).(string),
		RoleKey:   d.Get(projectRoleKey).(string),
	})
	if err != nil {
		return diag.Errorf("failed to delete project role: %v", err)
	}
	return nil
}

func updateProjectRole(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started update")

	clientinfo, ok := m.(*ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := getManagementClient(clientinfo, d.Get(projectRoleOrgID).(string))
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.UpdateProjectRole(ctx, &management2.UpdateProjectRoleRequest{
		ProjectId:   d.Get(projectRoleProjectID).(string),
		RoleKey:     d.Get(projectRoleKey).(string),
		DisplayName: d.Get(projectRoleDisplayName).(string),
		Group:       d.Get(projectRoleGroup).(string),
	})
	if err != nil {
		return diag.Errorf("failed to update project role: %v", err)
	}

	return nil
}

func createProjectRole(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started create")

	clientinfo, ok := m.(*ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	orgID := d.Get(projectRoleOrgID).(string)
	client, err := getManagementClient(clientinfo, orgID)
	if err != nil {
		return diag.FromErr(err)
	}

	projectID := d.Get(projectRoleProjectID).(string)
	roleKey := d.Get(projectRoleKey).(string)
	_, err = client.AddProjectRole(ctx, &management2.AddProjectRoleRequest{
		ProjectId:   projectID,
		RoleKey:     roleKey,
		DisplayName: d.Get(projectRoleDisplayName).(string),
		Group:       d.Get(projectRoleGroup).(string),
	})
	if err != nil {
		return diag.Errorf("failed to create project role: %v", err)
	}
	d.SetId(getProjectRoleID(orgID, projectID, roleKey))

	return nil
}

func readProjectRole(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started read")

	clientinfo, ok := m.(*ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := getManagementClient(clientinfo, d.Get(projectRoleOrgID).(string))
	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := client.ListProjectRoles(ctx, &management2.ListProjectRolesRequest{
		ProjectId: d.Get(projectRoleProjectID).(string),
		Queries: []*project2.RoleQuery{
			{Query: &project2.RoleQuery_KeyQuery{
				KeyQuery: &project2.RoleKeyQuery{
					Key:    d.Get(projectRoleKey).(string),
					Method: object.TextQueryMethod_TEXT_QUERY_METHOD_EQUALS,
				},
			}},
		},
	})
	if err != nil || resp.Result == nil || len(resp.Result) == 0 {
		d.SetId("")
		return nil
		//return diag.Errorf("failed to read project role: %v", err)
	}

	if len(resp.Result) == 1 {
		projectID := d.Get(projectRoleProjectID).(string)
		orgID := d.Get(projectRoleOrgID).(string)
		projectRole := resp.GetResult()[0]
		roleKey := projectRole.GetKey()
		set := map[string]interface{}{
			projectRoleProjectID:   projectID,
			projectRoleOrgID:       orgID,
			projectRoleKey:         roleKey,
			projectRoleDisplayName: projectRole.GetDisplayName(),
			projectRoleGroup:       projectRole.GetGroup(),
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
