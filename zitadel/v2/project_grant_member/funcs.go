package project_grant_member

import (
	"context"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/management"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/member"

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

	_, err = client.RemoveProjectGrantMember(ctx, &management.RemoveProjectGrantMemberRequest{
		UserId:    d.Get(userIDVar).(string),
		ProjectId: d.Get(projectIDVar).(string),
		GrantId:   d.Get(grantIDVar).(string),
	})
	if err != nil {
		return diag.Errorf("failed to delete projectmember: %v", err)
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

	_, err = client.UpdateProjectGrantMember(ctx, &management.UpdateProjectGrantMemberRequest{
		UserId:    d.Get(userIDVar).(string),
		Roles:     helper.GetOkSetToStringSlice(d, rolesVar),
		ProjectId: d.Get(projectIDVar).(string),
		GrantId:   d.Get(grantIDVar).(string),
	})
	if err != nil {
		return diag.Errorf("failed to update projectmember: %v", err)
	}
	return nil
}

func create(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started create")

	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	org := d.Get(orgIDVar).(string)
	client, err := helper.GetManagementClient(clientinfo, org)
	if err != nil {
		return diag.FromErr(err)
	}

	userID := d.Get(userIDVar).(string)
	projectID := d.Get(projectIDVar).(string)
	grantID := d.Get(grantIDVar).(string)
	_, err = client.AddProjectGrantMember(ctx, &management.AddProjectGrantMemberRequest{
		UserId:    userID,
		ProjectId: projectID,
		GrantId:   grantID,
		Roles:     helper.GetOkSetToStringSlice(d, rolesVar),
	})
	if err != nil {
		return diag.Errorf("failed to create projectgrantmember: %v", err)
	}
	d.SetId(getProjectGrantMemberID(org, projectID, grantID, userID))
	return nil
}

func read(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started read")

	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}
	org := d.Get(orgIDVar).(string)
	client, err := helper.GetManagementClient(clientinfo, org)
	if err != nil {
		return diag.FromErr(err)
	}

	projectID := d.Get(projectIDVar).(string)
	grantID := d.Get(grantIDVar).(string)
	userID := d.Get(userIDVar).(string)
	resp, err := client.ListProjectGrantMembers(ctx, &management.ListProjectGrantMembersRequest{
		ProjectId: projectID,
		GrantId:   grantID,
		Queries: []*member.SearchQuery{{
			Query: &member.SearchQuery_UserIdQuery{
				UserIdQuery: &member.UserIDQuery{
					UserId: userID,
				},
			},
		}},
	})
	if err != nil && helper.IgnoreIfNotFoundError(err) == nil {
		d.SetId("")
		return nil
	}
	if err != nil {
		return diag.Errorf("failed to list projectgrantmembers")
	}

	if len(resp.Result) == 1 {
		memberRes := resp.Result[0]
		set := map[string]interface{}{
			userIDVar:    userID,
			orgIDVar:     memberRes.GetDetails().GetResourceOwner(),
			projectIDVar: projectID,
			rolesVar:     memberRes.GetRoles(),
			grantIDVar:   grantID,
		}
		for k, v := range set {
			if err := d.Set(k, v); err != nil {
				return diag.Errorf("failed to set %s of projectgrantmember: %v", k, err)
			}
		}
		d.SetId(getProjectGrantMemberID(org, projectID, grantID, userID))
		return nil
	}

	d.SetId("")
	return nil
}

func getProjectGrantMemberID(org, projectID, grantID, userID string) string {
	return org + "_" + projectID + "_" + grantID + "_" + userID
}
