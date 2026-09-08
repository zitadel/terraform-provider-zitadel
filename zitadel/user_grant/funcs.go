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

func list(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started list")
	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}
	client, err := helper.GetManagementClient(ctx, clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}

	userID := d.Get(UserIDVar).(string)
	queries := []*user.UserGrantQuery{
		{
			Query: &user.UserGrantQuery_UserIdQuery{
				UserIdQuery: &user.UserGrantUserIDQuery{UserId: userID},
			},
		},
	}
	if projectID, ok := d.GetOk(projectIDVar); ok {
		queries = append(queries, &user.UserGrantQuery{
			Query: &user.UserGrantQuery_ProjectIdQuery{
				ProjectIdQuery: &user.UserGrantProjectIDQuery{ProjectId: projectID.(string)},
			},
		})
	}
	if projectGrantID, ok := d.GetOk(projectGrantIDVar); ok {
		queries = append(queries, &user.UserGrantQuery{
			Query: &user.UserGrantQuery_ProjectGrantIdQuery{
				ProjectGrantIdQuery: &user.UserGrantProjectGrantIDQuery{ProjectGrantId: projectGrantID.(string)},
			},
		})
	}
	if roleKey, ok := d.GetOk(roleKeyVar); ok {
		queries = append(queries, &user.UserGrantQuery{
			Query: &user.UserGrantQuery_RoleKeyQuery{
				RoleKeyQuery: &user.UserGrantRoleKeyQuery{RoleKey: roleKey.(string)},
			},
		})
	}

	resp, err := client.ListUserGrants(helper.CtxWithOrgID(ctx, d), &management.ListUserGrantRequest{
		Queries: queries,
	})
	if err != nil {
		return diag.Errorf("failed to list user grants: %v", err)
	}

	grants := make([]interface{}, len(resp.GetResult()))
	for i, grant := range resp.GetResult() {
		grants[i] = map[string]interface{}{
			idVar:             grant.GetId(),
			projectIDVar:      grant.GetProjectId(),
			projectNameVar:    grant.GetProjectName(),
			projectGrantIDVar: grant.GetProjectGrantId(),
			grantedOrgIDVar:   grant.GetGrantedOrgId(),
			RoleKeysVar:       grant.GetRoleKeys(),
			stateVar:          grant.GetState().String(),
		}
	}

	d.SetId(userID)
	return diag.FromErr(d.Set(userGrantsVar, grants))
}
