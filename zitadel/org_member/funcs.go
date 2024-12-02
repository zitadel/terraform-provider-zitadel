package org_member

import (
	"context"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/management"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/member"

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

	_, err = client.RemoveOrgMember(helper.CtxWithOrgID(ctx, d), &management.RemoveOrgMemberRequest{
		UserId: d.Get(UserIDVar).(string),
	})
	if err != nil {
		return diag.Errorf("failed to delete orgmember: %v", err)
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

	_, err = client.UpdateOrgMember(helper.CtxWithOrgID(ctx, d), &management.UpdateOrgMemberRequest{
		UserId: d.Get(UserIDVar).(string),
		Roles:  helper.GetOkSetToStringSlice(d, RolesVar),
	})
	if err != nil {
		return diag.Errorf("failed to update orgmember: %v", err)
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
	_, err = client.AddOrgMember(helper.CtxWithOrgID(ctx, d), &management.AddOrgMemberRequest{
		UserId: userID,
		Roles:  helper.GetOkSetToStringSlice(d, RolesVar),
	})
	if err != nil {
		return diag.Errorf("failed to create orgmember: %v", err)
	}
	d.SetId(getOrgMemberID(org, userID))
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

	userID := d.Get(UserIDVar).(string)
	resp, err := client.ListOrgMembers(helper.CtxWithOrgID(ctx, d), &management.ListOrgMembersRequest{
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
		return diag.Errorf("failed to list org members")
	}

	if len(resp.Result) == 1 {
		orgMember := resp.Result[0]
		set := map[string]interface{}{
			UserIDVar:       orgMember.GetUserId(),
			helper.OrgIDVar: orgMember.GetDetails().GetResourceOwner(),
			RolesVar:        orgMember.GetRoles(),
		}
		for k, v := range set {
			if err := d.Set(k, v); err != nil {
				return diag.Errorf("failed to set %s of orgmember: %v", k, err)
			}
		}
		d.SetId(getOrgMemberID(org, userID))
		return nil
	}

	d.SetId("")
	return nil
}

func getOrgMemberID(org string, userID string) string {
	return org + "_" + userID
}
