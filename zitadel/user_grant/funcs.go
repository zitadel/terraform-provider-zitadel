package user_grant

import (
	"context"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/management"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/user"

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

	_, err = client.RemoveUserGrant(helper.CtxWithOrgID(ctx, d), &management.RemoveUserGrantRequest{
		GrantId: d.Id(),
		UserId:  d.Get(UserIDVar).(string),
	})
	if err != nil {
		return diag.Errorf("failed to delete usergrant: %v", err)
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

	_, err = client.UpdateUserGrant(helper.CtxWithOrgID(ctx, d), &management.UpdateUserGrantRequest{
		GrantId:  d.Id(),
		UserId:   d.Get(UserIDVar).(string),
		RoleKeys: helper.GetOkSetToStringSlice(d, RoleKeysVar),
	})
	if err != nil {
		return diag.Errorf("failed to update usergrant: %v", err)
	}
	return nil
}

func create(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started create")

	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := helper.GetManagementClient(ctx, clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := client.AddUserGrant(helper.CtxWithOrgID(ctx, d), &management.AddUserGrantRequest{
		UserId:         d.Get(UserIDVar).(string),
		ProjectGrantId: d.Get(projectGrantIDVar).(string),
		ProjectId:      d.Get(projectIDVar).(string),
		RoleKeys:       helper.GetOkSetToStringSlice(d, RoleKeysVar),
	})
	if err != nil {
		return diag.Errorf("failed to create usergrant: %v", err)
	}
	d.SetId(resp.GetUserGrantId())
	return nil
}

func read(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started read")
	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}
	client, err := helper.GetManagementClient(ctx, clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}
	resp, err := client.GetUserGrantByID(helper.CtxWithOrgID(ctx, d), &management.GetUserGrantByIDRequest{
		GrantId: helper.GetID(d, grantIDVar),
		UserId:  d.Get(UserIDVar).(string),
	})
	if err != nil && helper.IgnoreIfNotFoundError(err) == nil {
		d.SetId("")
		return nil
	}
	if err != nil {
		return diag.Errorf("failed to get user grant: %v", err)
	}
	grant := resp.GetUserGrant()
	set := map[string]interface{}{
		UserIDVar:       grant.GetUserId(),
		RoleKeysVar:     grant.GetRoleKeys(),
		helper.OrgIDVar: grant.GetDetails().GetResourceOwner(),
	}
	if grant.GetProjectId() != "" {
		set[projectIDVar] = grant.GetProjectId()
	}
	if grant.GetProjectGrantId() != "" {
		set[projectGrantIDVar] = grant.GetProjectGrantId()
	}
	for k, v := range set {
		if err := d.Set(k, v); err != nil {
			return diag.Errorf("failed to set %s of usergrant: %v", k, err)
		}
	}
	d.SetId(grant.GetId())
	return nil
}

func get(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started get")
	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}
	client, err := helper.GetManagementClient(ctx, clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}

	grantID := helper.GetID(d, grantIDVar)
	userID := helper.GetID(d, UserIDVar)

	resp, err := client.GetUserGrantByID(helper.CtxWithOrgID(ctx, d), &management.GetUserGrantByIDRequest{GrantId: grantID, UserId: userID})
	if err != nil && helper.IgnoreIfNotFoundError(err) == nil {
		d.SetId("")
		return nil
	}
	if err != nil {
		return diag.Errorf("failed to get user grant: %v", err)
	}

	grant := resp.GetUserGrant()
	d.SetId(grant.GetId())
	if err := d.Set(UserIDVar, grant.GetUserId()); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set(projectIDVar, grant.GetProjectId()); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set(projectGrantIDVar, grant.GetProjectGrantId()); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set(RoleKeysVar, grant.GetRoleKeys()); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func list(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started list")
	userID := d.Get(UserIDVar).(string)
	projectID := d.Get(projectIDVar).(string)
	projectGrantID := d.Get(projectGrantIDVar).(string)
	roleKey := d.Get(roleKeyVar).(string)
	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}
	client, err := helper.GetManagementClient(ctx, clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}
	req := &management.ListUserGrantRequest{}
	req.Queries = append(req.Queries, &user.UserGrantQuery{
		Query: &user.UserGrantQuery_UserIdQuery{
			UserIdQuery: &user.UserGrantUserIDQuery{
				UserId: userID,
			},
		},
	})
	if projectID != "" {
		req.Queries = append(req.Queries, &user.UserGrantQuery{
			Query: &user.UserGrantQuery_ProjectIdQuery{
				ProjectIdQuery: &user.UserGrantProjectIDQuery{
					ProjectId: projectID,
				},
			},
		})
	}
	if projectGrantID != "" {
		req.Queries = append(req.Queries, &user.UserGrantQuery{
			Query: &user.UserGrantQuery_ProjectGrantIdQuery{
				ProjectGrantIdQuery: &user.UserGrantProjectGrantIDQuery{
					ProjectGrantId: projectGrantID,
				},
			},
		})
	}
	if roleKey != "" {
		req.Queries = append(req.Queries, &user.UserGrantQuery{
			Query: &user.UserGrantQuery_RoleKeyQuery{
				RoleKeyQuery: &user.UserGrantRoleKeyQuery{
					RoleKey: roleKey,
				},
			},
		})
	}
	resp, err := client.ListUserGrants(helper.CtxWithOrgID(ctx, d), req)
	if err != nil {
		return diag.Errorf("error while getting user grant list: %v", err)
	}
	grantIDs := make([]string, 0, len(resp.Result))
	for _, grant := range resp.Result {
		grantIDs = append(grantIDs, grant.Id)
	}
	d.SetId("-")
	return diag.FromErr(d.Set(grantIDsVar, grantIDs))
}
