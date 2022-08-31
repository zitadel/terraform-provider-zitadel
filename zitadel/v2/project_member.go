package v2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	management2 "github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/management"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/member"
)

const (
	projectMemberOrgIDVar     = "org_id"
	projectMemberProjectIDVar = "project_id"
	projectMemberUserIDVar    = "user_id"
	projectMemberRolesVar     = "roles"
)

func GetProjectMember() *schema.Resource {
	return &schema.Resource{
		Description: "Resource representing the membership of a user on an project, defined with the given role.",
		Schema: map[string]*schema.Schema{
			projectMemberOrgIDVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the organization which owns the resource",
				ForceNew:    true,
			},
			projectMemberProjectIDVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the project",
				ForceNew:    true,
			},
			projectMemberUserIDVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the user",
				ForceNew:    true,
			},
			projectMemberRolesVar: {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Required:    true,
				Description: "List of roles granted",
			},
		},
		DeleteContext: deleteProjectMember,
		CreateContext: createProjectMember,
		UpdateContext: updateProjectMember,
		ReadContext:   readProjectMember,
		Importer:      &schema.ResourceImporter{StateContext: schema.ImportStatePassthroughContext},
	}
}

func deleteProjectMember(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started delete")

	clientinfo, ok := m.(*ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := getManagementClient(clientinfo, d.Get(projectMemberOrgIDVar).(string))
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.RemoveProjectMember(ctx, &management2.RemoveProjectMemberRequest{
		UserId:    d.Get(projectMemberUserIDVar).(string),
		ProjectId: d.Get(projectMemberProjectIDVar).(string),
	})
	if err != nil {
		return diag.Errorf("failed to delete projectmember: %v", err)
	}
	return nil
}

func updateProjectMember(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started update")

	clientinfo, ok := m.(*ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := getManagementClient(clientinfo, d.Get(projectMemberOrgIDVar).(string))
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.UpdateProjectMember(ctx, &management2.UpdateProjectMemberRequest{
		UserId:    d.Get(projectMemberUserIDVar).(string),
		Roles:     d.Get(projectMemberRolesVar).([]string),
		ProjectId: d.Get(projectMemberProjectIDVar).(string),
	})
	if err != nil {
		return diag.Errorf("failed to update projectmember: %v", err)
	}
	return nil
}

func createProjectMember(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started create")

	clientinfo, ok := m.(*ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	org := d.Get(projectMemberOrgIDVar).(string)
	client, err := getManagementClient(clientinfo, org)
	if err != nil {
		return diag.FromErr(err)
	}

	userID := d.Get(projectMemberUserIDVar).(string)
	projectID := d.Get(projectMemberProjectIDVar).(string)
	roles := make([]string, 0)
	for _, role := range d.Get(projectMemberRolesVar).(*schema.Set).List() {
		roles = append(roles, role.(string))
	}

	_, err = client.AddProjectMember(ctx, &management2.AddProjectMemberRequest{
		UserId:    userID,
		ProjectId: projectID,
		Roles:     roles,
	})
	if err != nil {
		return diag.Errorf("failed to create projectmember: %v", err)
	}
	d.SetId(getProjectMemberID(org, projectID, userID))
	return nil
}

func readProjectMember(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started read")

	clientinfo, ok := m.(*ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}
	org := d.Get(projectMemberOrgIDVar).(string)
	client, err := getManagementClient(clientinfo, org)
	if err != nil {
		return diag.FromErr(err)
	}

	projectID := d.Get(projectMemberProjectIDVar).(string)
	userID := d.Get(projectMemberUserIDVar).(string)
	resp, err := client.ListProjectMembers(ctx, &management2.ListProjectMembersRequest{
		ProjectId: projectID,
		Queries: []*member.SearchQuery{{
			Query: &member.SearchQuery_UserIdQuery{
				UserIdQuery: &member.UserIDQuery{
					UserId: userID,
				},
			},
		}},
	})
	if err != nil {
		d.SetId("")
		return nil
		//return diag.Errorf("failed to read projectmember: %v", err)
	}

	if len(resp.Result) == 1 {
		memberRes := resp.Result[0]
		set := map[string]interface{}{
			projectMemberUserIDVar:    memberRes.GetUserId(),
			projectMemberOrgIDVar:     memberRes.GetDetails().GetResourceOwner(),
			projectMemberProjectIDVar: projectID,
			projectMemberRolesVar:     memberRes.GetRoles(),
		}
		for k, v := range set {
			if err := d.Set(k, v); err != nil {
				return diag.Errorf("failed to set %s of projectmember: %v", k, err)
			}
		}
		d.SetId(getProjectMemberID(org, projectID, userID))
		return nil
	}

	d.SetId("")
	return nil
}

func getProjectMemberID(org string, projectID string, userID string) string {
	return org + "_" + projectID + "_" + userID
}
