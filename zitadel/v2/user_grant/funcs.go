package user_grant

import (
	"context"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/management"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/user"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/helper"
)

func delete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started delete")

	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := helper.GetManagementClient(clientinfo, d.Get(userGrantOrgIDVar).(string))
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.RemoveUserGrant(ctx, &management.RemoveUserGrantRequest{
		GrantId: d.Id(),
		UserId:  d.Get(userGrantUserIDVar).(string),
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

	client, err := helper.GetManagementClient(clientinfo, d.Get(userGrantOrgIDVar).(string))
	if err != nil {
		return diag.FromErr(err)
	}

	roles := make([]string, 0)
	for _, role := range d.Get(userGrantRoleKeysVar).(*schema.Set).List() {
		roles = append(roles, role.(string))
	}
	_, err = client.UpdateUserGrant(ctx, &management.UpdateUserGrantRequest{
		GrantId:  d.Id(),
		UserId:   d.Get(userGrantUserIDVar).(string),
		RoleKeys: roles,
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

	client, err := helper.GetManagementClient(clientinfo, d.Get(userGrantOrgIDVar).(string))
	if err != nil {
		return diag.FromErr(err)
	}

	roles := make([]string, 0)
	for _, role := range d.Get(userGrantRoleKeysVar).(*schema.Set).List() {
		roles = append(roles, role.(string))
	}

	resp, err := client.AddUserGrant(ctx, &management.AddUserGrantRequest{
		UserId:         d.Get(userGrantUserIDVar).(string),
		ProjectGrantId: d.Get(userGrantProjectGrantIDVar).(string),
		ProjectId:      d.Get(userGrantProjectIDVar).(string),
		RoleKeys:       roles,
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

	client, err := helper.GetManagementClient(clientinfo, d.Get(userGrantOrgIDVar).(string))
	if err != nil {
		return diag.FromErr(err)
	}
	projectID := d.Get(userGrantProjectIDVar)
	projectGrantID := d.Get(userGrantProjectGrantIDVar)

	queries := []*user.UserGrantQuery{
		{Query: &user.UserGrantQuery_UserIdQuery{
			UserIdQuery: &user.UserGrantUserIDQuery{
				UserId: d.Get(userGrantUserIDVar).(string),
			},
		}},
	}
	if projectID != nil {
		queries = append(queries, &user.UserGrantQuery{Query: &user.UserGrantQuery_ProjectIdQuery{
			ProjectIdQuery: &user.UserGrantProjectIDQuery{
				ProjectId: projectID.(string),
			},
		}},
		)
	}
	if projectGrantID != nil {
		queries = append(queries, &user.UserGrantQuery{Query: &user.UserGrantQuery_ProjectGrantIdQuery{
			ProjectGrantIdQuery: &user.UserGrantProjectGrantIDQuery{
				ProjectGrantId: projectGrantID.(string),
			},
		}},
		)
	}
	grants, err := client.ListUserGrants(ctx, &management.ListUserGrantRequest{
		Queries: queries,
	})
	if err != nil {
		d.SetId("")
		return nil
		//return diag.Errorf("failed to read usergrant: %v", err)
	}

	if len(grants.GetResult()) == 1 {
		grant := grants.GetResult()[0]
		set := map[string]interface{}{
			userGrantUserIDVar:   grant.GetUserId(),
			userGrantRoleKeysVar: grant.GetRoleKeys(),
			userGrantOrgIDVar:    grant.GetDetails().GetResourceOwner(),
		}
		if grant.GetProjectId() != "" {
			set[userGrantProjectIDVar] = grant.GetProjectId()
		}
		if grant.GetProjectGrantId() != "" {
			set[userGrantProjectGrantIDVar] = grant.GetProjectGrantId()
		}
		for k, v := range set {
			if err := d.Set(k, v); err != nil {
				return diag.Errorf("failed to set %s of usergrant: %v", k, err)
			}
		}
		d.SetId(grant.GetId())
		return nil
	}

	d.SetId("")
	return nil
}
