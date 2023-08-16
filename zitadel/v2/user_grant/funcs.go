package user_grant

import (
	"context"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/management"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/helper"
)

func delete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started delete")

	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := helper.GetManagementClient(clientinfo, d.Get(helper.OrgIDVar).(string))
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.RemoveUserGrant(ctx, &management.RemoveUserGrantRequest{
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

	client, err := helper.GetManagementClient(clientinfo, d.Get(helper.OrgIDVar).(string))
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.UpdateUserGrant(ctx, &management.UpdateUserGrantRequest{
		GrantId:  d.Id(),
		UserId:   d.Get(UserIDVar).(string),
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

	client, err := helper.GetManagementClient(clientinfo, d.Get(helper.OrgIDVar).(string))
	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := client.AddUserGrant(ctx, &management.AddUserGrantRequest{
		UserId:         d.Get(UserIDVar).(string),
		ProjectGrantId: d.Get(ProjectGrantIDVar).(string),
		ProjectId:      d.Get(ProjectIDVar).(string),
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

	client, err := helper.GetManagementClient(clientinfo, d.Get(helper.OrgIDVar).(string))
	if err != nil {
		return diag.FromErr(err)
	}
	resp, err := client.GetUserGrantByID(ctx, &management.GetUserGrantByIDRequest{
		UserId:  d.Get(UserIDVar).(string),
		GrantId: helper.GetID(d, helper.ResourceIDVar),
	})
	if err != nil {
		return diag.Errorf("failed to get usergrants")
	}
	grant := resp.GetUserGrant()
	set := map[string]interface{}{
		UserIDVar:       grant.GetUserId(),
		roleKeysVar:     grant.GetRoleKeys(),
		helper.OrgIDVar: grant.GetDetails().GetResourceOwner(),
	}
	if grant.GetProjectId() != "" {
		set[ProjectIDVar] = grant.GetProjectId()
	}
	if grant.GetProjectGrantId() != "" {
		set[ProjectGrantIDVar] = grant.GetProjectGrantId()
	}
	for k, v := range set {
		if err := d.Set(k, v); err != nil {
			return diag.Errorf("failed to set %s of usergrant: %v", k, err)
		}
	}
	d.SetId(grant.GetId())
	return nil
}
