package project_member

import (
	"context"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/management"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/member"

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

	_, err = client.RemoveProjectMember(helper.CtxWithOrgID(ctx, d), &management.RemoveProjectMemberRequest{
		UserId:    d.Get(UserIDVar).(string),
		ProjectId: d.Get(ProjectIDVar).(string),
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

	client, err := helper.GetManagementClient(clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.UpdateProjectMember(helper.CtxWithOrgID(ctx, d), &management.UpdateProjectMemberRequest{
		UserId:    d.Get(UserIDVar).(string),
		Roles:     helper.GetOkSetToStringSlice(d, rolesVar),
		ProjectId: d.Get(ProjectIDVar).(string),
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
	client, err := helper.GetManagementClient(clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}

	userID := d.Get(UserIDVar).(string)
	projectID := d.Get(ProjectIDVar).(string)
	_, err = client.AddProjectMember(helper.CtxWithOrgID(ctx, d), &management.AddProjectMemberRequest{
		UserId:    userID,
		ProjectId: projectID,
		Roles:     helper.GetOkSetToStringSlice(d, rolesVar),
	})
	if err != nil {
		return diag.Errorf("failed to create projectmember: %v", err)
	}
	d.SetId(getProjectMemberID(org, projectID, userID))
	return nil
}

func read(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started read")

	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}
	org := d.Get(helper.OrgIDVar).(string)
	client, err := helper.GetManagementClient(clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}

	projectID := d.Get(ProjectIDVar).(string)
	userID := d.Get(UserIDVar).(string)
	resp, err := client.ListProjectMembers(helper.CtxWithOrgID(ctx, d), &management.ListProjectMembersRequest{
		ProjectId: projectID,
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
		return diag.Errorf("failed to list projectmembers")
	}

	if len(resp.Result) == 1 {
		memberRes := resp.Result[0]
		set := map[string]interface{}{
			UserIDVar:       memberRes.GetUserId(),
			helper.OrgIDVar: memberRes.GetDetails().GetResourceOwner(),
			ProjectIDVar:    projectID,
			rolesVar:        memberRes.GetRoles(),
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
