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

	client, err := helper.GetManagementClient(clientinfo, d.Get(orgIDVar).(string))
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.RemoveUserGrant(ctx, &management.RemoveUserGrantRequest{
		GrantId: d.Id(),
		UserId:  d.Get(userIDVar).(string),
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

	client, err := helper.GetManagementClient(clientinfo, d.Get(orgIDVar).(string))
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.UpdateUserGrant(ctx, &management.UpdateUserGrantRequest{
		GrantId:  d.Id(),
		UserId:   d.Get(userIDVar).(string),
		RoleKeys: helper.GetOkSetToStringSlice(d, roleKeysVar),
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

	client, err := helper.GetManagementClient(clientinfo, d.Get(orgIDVar).(string))
	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := client.AddUserGrant(ctx, &management.AddUserGrantRequest{
		UserId:         d.Get(userIDVar).(string),
		ProjectGrantId: d.Get(projectGrantIDVar).(string),
		ProjectId:      d.Get(projectIDVar).(string),
		RoleKeys:       helper.GetOkSetToStringSlice(d, roleKeysVar),
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

	client, err := helper.GetManagementClient(clientinfo, d.Get(orgIDVar).(string))
	if err != nil {
		return diag.FromErr(err)
	}
	projectID := d.Get(projectIDVar)
	projectGrantID := d.Get(projectGrantIDVar)

	queries := []*user.UserGrantQuery{
		{Query: &user.UserGrantQuery_UserIdQuery{
			UserIdQuery: &user.UserGrantUserIDQuery{
				UserId: d.Get(userIDVar).(string),
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
		return diag.Errorf("failed to list usergrants")
	}

	if len(grants.GetResult()) == 1 {
		grant := grants.GetResult()[0]
		set := map[string]interface{}{
			userIDVar:   grant.GetUserId(),
			roleKeysVar: grant.GetRoleKeys(),
			orgIDVar:    grant.GetDetails().GetResourceOwner(),
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

	d.SetId("")
	return nil
}
