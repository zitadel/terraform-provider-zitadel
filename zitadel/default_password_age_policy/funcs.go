package default_password_age_policy

import (
	"context"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/admin"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
)

func delete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "default password age policy cannot be deleted")
	return nil
}

func update(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started update")

	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := helper.GetAdminClient(ctx, clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}

	id := ""
	if d.HasChanges(maxAgeDays, expireWarnDays) {
		req := admin.UpdatePasswordAgePolicyRequest{
			MaxAgeDays:     uint32(d.Get(maxAgeDays).(int)),
			ExpireWarnDays: uint32(d.Get(expireWarnDays).(int)),
		}

		resp, err := client.UpdatePasswordAgePolicy(ctx, &req)
		if err != nil {
			return diag.Errorf("failed to update default password age policy: %v", err)
		}

		if helper.IgnorePreconditionError(err) != nil {
			return diag.Errorf("failed to update default password age policy: %v", err)
		}

		if resp != nil {
			id = resp.GetDetails().GetResourceOwner()
		}
	}

	if id == "" {
		resp, err := client.GetPasswordAgePolicy(ctx, &admin.GetPasswordAgePolicyRequest{})
		if err != nil {
			return diag.Errorf("failed to get default password complexity policy: %v", err)
		}

		id = resp.GetPolicy().GetDetails().GetResourceOwner()
	}
	d.SetId(id)

	return nil
}

func read(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started read")

	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := helper.GetAdminClient(ctx, clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := client.GetPasswordAgePolicy(ctx, &admin.GetPasswordAgePolicyRequest{})
	if err != nil && helper.IgnoreIfNotFoundError(err) == nil {
		d.SetId("")
		return nil
	}

	if err != nil {
		return diag.Errorf("failed to get default password age policy: %v", err)
	}

	policy := resp.Policy
	set := map[string]interface{}{
		maxAgeDays:     policy.GetMaxAgeDays(),
		expireWarnDays: policy.GetExpireWarnDays(),
	}

	for k, v := range set {
		if err := d.Set(k, v); err != nil {
			return diag.Errorf("failed to set %s of password age policy: %v", k, err)
		}
	}

	d.SetId(policy.GetDetails().GetResourceOwner())
	return nil
}
