package org

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/admin"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/management"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/helper"
)

func delete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "org delete not yet implemented")
	return nil
}

func create(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started create")

	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := helper.GetManagementClient(clientinfo, "")
	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := client.AddOrg(ctx, &management.AddOrgRequest{
		Name: d.Get(nameVar).(string),
	})
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(resp.GetId())

	return nil
}

func update(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started update")

	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := helper.GetManagementClient(clientinfo, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.UpdateOrg(ctx, &management.UpdateOrgRequest{
		Name: d.Get(nameVar).(string),
	})
	if err != nil {
		return diag.Errorf("failed to update org: %v", err)
	}
	return nil
}

func read(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started read")

	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := helper.GetAdminClient(clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := client.ListOrgs(ctx, &admin.ListOrgsRequest{})
	if err != nil {
		return diag.Errorf("error while listing orgs: %v", err)
	}
	tflog.Debug(ctx, "found orgs", map[string]interface{}{
		"orglist": resp.Result,
	})

	orgID := helper.GetID(d, orgIDVar)
	tflog.Debug(ctx, "check if org is existing", map[string]interface{}{
		"id": orgID,
	})

	for i := range resp.Result {
		org := resp.Result[i]
		if strings.Compare(org.GetId(), orgID) == 0 {
			d.SetId(orgID)
			tflog.Debug(ctx, "found org", map[string]interface{}{
				"id": orgID,
			})
			if err := d.Set(nameVar, org.GetName()); err != nil {
				return diag.Errorf("failed to set %s of org: %v", nameVar, err)
			}
			return nil
		}
	}

	d.SetId("")
	tflog.Debug(ctx, "org not found", map[string]interface{}{})
	return nil
}
