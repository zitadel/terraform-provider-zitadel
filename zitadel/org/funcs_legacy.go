package org

import (
	"context"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/admin"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/management"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
)

func legacyDeleteOrg(ctx context.Context, d *schema.ResourceData, clientinfo *helper.ClientInfo) diag.Diagnostics {
	tflog.Info(ctx, "falling back to legacy admin API for delete")

	client, err := helper.GetAdminClient(ctx, clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.RemoveOrg(ctx, &admin.RemoveOrgRequest{
		OrgId: d.Id(),
	})
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return nil
}

func legacyUpdateOrg(ctx context.Context, d *schema.ResourceData, clientinfo *helper.ClientInfo) diag.Diagnostics {
	tflog.Info(ctx, "falling back to legacy management API for update")

	client, err := helper.GetManagementClient(ctx, clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.UpdateOrg(helper.CtxSetOrgID(ctx, d.Id()), &management.UpdateOrgRequest{
		Name: d.Get(NameVar).(string),
	})
	if err != nil {
		return diag.Errorf("failed to update org: %v", err)
	}
	return nil
}
