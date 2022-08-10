package v2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	management2 "github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/management"
)

const (
	projectGrantMemberOrgIDVar     = "org_id"
	projectGrantMemberProjectIDVar = "project_id"
	projectGrantMemberGrantIDVar   = "grant_id"
	projectGrantMemberUserIDVar    = "user_id"
	projectGrantMemberRolesVar     = "roles"
)

func GetProjectGrantMember() *schema.Resource {
	return &schema.Resource{
		Description: "Resource representing the membership of a user on an granted project, defined with the given role.",
		Schema: map[string]*schema.Schema{
			projectGrantMemberOrgIDVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the organization which owns the resource",
				ForceNew:    true,
			},
			projectGrantMemberProjectIDVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the project",
				ForceNew:    true,
			},
			projectGrantMemberGrantIDVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the grant",
				ForceNew:    true,
			},
			projectGrantMemberUserIDVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the user",
				ForceNew:    true,
			},
			projectGrantMemberRolesVar: {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Required:    true,
				Description: "List of roles granted",
			},
		},
		DeleteContext: deleteProjectGrantMember,
		CreateContext: createProjectGrantMember,
		UpdateContext: updateProjectGrantMember,
		ReadContext:   readProjectGrantMember,
	}
}

func deleteProjectGrantMember(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started delete")

	clientinfo, ok := m.(*ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := getManagementClient(clientinfo, d.Get(projectGrantMemberOrgIDVar).(string))
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.RemoveProjectGrantMember(ctx, &management2.RemoveProjectGrantMemberRequest{
		UserId:    d.Get(projectGrantMemberUserIDVar).(string),
		ProjectId: d.Get(projectGrantMemberProjectIDVar).(string),
		GrantId:   d.Get(projectGrantMemberGrantIDVar).(string),
	})
	if err != nil {
		return diag.Errorf("failed to delete projectmember: %v", err)
	}
	return nil
}

func updateProjectGrantMember(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started update")

	clientinfo, ok := m.(*ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := getManagementClient(clientinfo, d.Get(projectGrantMemberOrgIDVar).(string))
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.UpdateProjectGrantMember(ctx, &management2.UpdateProjectGrantMemberRequest{
		UserId:    d.Get(projectGrantMemberUserIDVar).(string),
		Roles:     d.Get(projectGrantMemberRolesVar).([]string),
		ProjectId: d.Get(projectGrantMemberProjectIDVar).(string),
		GrantId:   d.Get(projectGrantMemberGrantIDVar).(string),
	})
	if err != nil {
		return diag.Errorf("failed to update projectmember: %v", err)
	}
	return nil
}

func createProjectGrantMember(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started create")

	clientinfo, ok := m.(*ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	org := d.Get(projectGrantMemberOrgIDVar).(string)
	client, err := getManagementClient(clientinfo, org)
	if err != nil {
		return diag.FromErr(err)
	}

	userID := d.Get(projectGrantMemberUserIDVar).(string)
	projectID := d.Get(projectGrantMemberProjectIDVar).(string)
	grantID := d.Get(projectGrantMemberGrantIDVar).(string)
	roles := make([]string, 0)
	for _, role := range d.Get(projectGrantMemberRolesVar).(*schema.Set).List() {
		roles = append(roles, role.(string))
	}
	_, err = client.AddProjectGrantMember(ctx, &management2.AddProjectGrantMemberRequest{
		UserId:    userID,
		ProjectId: projectID,
		GrantId:   grantID,
		Roles:     roles,
	})
	if err != nil {
		return diag.Errorf("failed to create projectgrantmember: %v", err)
	}
	d.SetId(getProjectGrantMemberID(org, projectID, grantID, userID))
	return nil
}

func readProjectGrantMember(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started read")

	clientinfo, ok := m.(*ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}
	org := d.Get(projectGrantMemberOrgIDVar).(string)
	client, err := getManagementClient(clientinfo, org)
	if err != nil {
		return diag.FromErr(err)
	}

	projectID := d.Get(projectGrantMemberProjectIDVar).(string)
	grantID := d.Get(projectGrantMemberGrantIDVar).(string)
	resp, err := client.ListProjectGrantMembers(ctx, &management2.ListProjectGrantMembersRequest{ProjectId: projectID, GrantId: grantID})
	if err != nil {
		d.SetId("")
		return nil
		//return diag.Errorf("failed to read projectgrantmember: %v", err)
	}

	userID := d.Get(projectGrantMemberUserIDVar).(string)
	for _, member := range resp.Result {
		if member.UserId == userID {
			set := map[string]interface{}{
				projectGrantMemberUserIDVar:    member.GetUserId(),
				projectGrantMemberOrgIDVar:     member.GetDetails().GetResourceOwner(),
				projectGrantMemberProjectIDVar: projectID,
				projectGrantMemberRolesVar:     member.GetRoles(),
				projectGrantMemberGrantIDVar:   grantID,
			}
			for k, v := range set {
				if err := d.Set(k, v); err != nil {
					return diag.Errorf("failed to set %s of projectgrantmember: %v", k, err)
				}
			}
			d.SetId(getProjectGrantMemberID(org, projectID, grantID, userID))
			return nil
		}
	}
	d.SetId("")
	return nil
}

func getProjectGrantMemberID(org, projectID, grantID, userID string) string {
	return org + "_" + projectID + "_" + grantID + "_" + userID
}
