package project_grant_member

import (
	"context"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/management"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/member"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/helper"
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

	_, err = client.RemoveProjectGrantMember(helper.CtxWithOrgID(ctx, d), &management.RemoveProjectGrantMemberRequest{
		UserId:    d.Get(UserIDVar).(string),
		ProjectId: d.Get(ProjectIDVar).(string),
		GrantId:   d.Get(GrantIDVar).(string),
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

	client, err := helper.GetManagementClient(ctx, clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.UpdateProjectGrantMember(helper.CtxWithOrgID(ctx, d), &management.UpdateProjectGrantMemberRequest{
		UserId:    d.Get(UserIDVar).(string),
		Roles:     helper.GetOkSetToStringSlice(d, RolesVar),
		ProjectId: d.Get(ProjectIDVar).(string),
		GrantId:   d.Get(GrantIDVar).(string),
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

	org := d.Get(helper.OrgIDVar).(string)
	client, err := helper.GetManagementClient(ctx, clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}

	userID := d.Get(UserIDVar).(string)
	projectID := d.Get(ProjectIDVar).(string)
	grantID := d.Get(GrantIDVar).(string)
	_, err = client.AddProjectGrantMember(helper.CtxWithOrgID(ctx, d), &management.AddProjectGrantMemberRequest{
		UserId:    userID,
		ProjectId: projectID,
		GrantId:   grantID,
		Roles:     helper.GetOkSetToStringSlice(d, RolesVar),
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
	org := d.Get(helper.OrgIDVar).(string)
	client, err := helper.GetManagementClient(ctx, clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}

	projectID := d.Get(ProjectIDVar).(string)
	grantID := d.Get(GrantIDVar).(string)
	userID := d.Get(UserIDVar).(string)
	resp, err := client.ListProjectGrantMembers(helper.CtxWithOrgID(ctx, d), &management.ListProjectGrantMembersRequest{
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
			UserIDVar:       userID,
			helper.OrgIDVar: org,
			ProjectIDVar:    projectID,
			RolesVar:        memberRes.GetRoles(),
			GrantIDVar:      grantID,
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
