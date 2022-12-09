package project_grant

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

	client, err := helper.GetManagementClient(clientinfo, d.Get(orgIDVar).(string))
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.RemoveProjectGrant(ctx, &management.RemoveProjectGrantRequest{
		GrantId:   d.Id(),
		ProjectId: d.Get(projectIDVar).(string),
	})
	if err != nil {
		return diag.Errorf("failed to delete projectgrant: %v", err)
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

	_, err = client.UpdateProjectGrant(ctx, &management.UpdateProjectGrantRequest{
		GrantId:   d.Id(),
		ProjectId: d.Get(projectIDVar).(string),
		RoleKeys:  helper.GetOkSetToStringSlice(d, roleKeysVar),
	})
	if err != nil {
		return diag.Errorf("failed to update projectgrant: %v", err)
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

	resp, err := client.AddProjectGrant(ctx, &management.AddProjectGrantRequest{
		GrantedOrgId: d.Get(grantedOrgIDVar).(string),
		ProjectId:    d.Get(projectIDVar).(string),
		RoleKeys:     helper.GetOkSetToStringSlice(d, roleKeysVar),
	})
	if err != nil {
		return diag.Errorf("failed to create projectgrant: %v", err)
	}
	d.SetId(resp.GetGrantId())
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

	resp, err := client.GetProjectGrantByID(ctx, &management.GetProjectGrantByIDRequest{ProjectId: d.Get(projectIDVar).(string), GrantId: d.Id()})
	if err != nil && helper.IgnoreIfNotFoundError(err) == nil {
		d.SetId("")
		return nil
	}
	if err != nil {
		return diag.Errorf("failed to get projectgrant")
	}

	projectGrant := resp.GetProjectGrant()
	set := map[string]interface{}{
		projectIDVar:    projectGrant.GetProjectId(),
		grantedOrgIDVar: projectGrant.GetGrantedOrgId(),
		roleKeysVar:     projectGrant.GetGrantedRoleKeys(),
		orgIDVar:        projectGrant.GetDetails().GetResourceOwner(),
	}
	for k, v := range set {
		if err := d.Set(k, v); err != nil {
			return diag.Errorf("failed to set %s of projectgrant: %v", k, err)
		}
	}
	d.SetId(projectGrant.GetGrantId())
	return nil
}
