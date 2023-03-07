package default_lockout_policy

import (
	"context"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/admin"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/helper"
)

func delete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "default lockout policy cannot be deleted")
	return nil
}

func update(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started update")

	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := helper.GetAdminClient(clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}

	id := ""
	if d.HasChanges(maxPasswordAttemptsVar) {
		resp, err := client.UpdateLockoutPolicy(ctx, &admin.UpdateLockoutPolicyRequest{
			MaxPasswordAttempts: uint32(d.Get(maxPasswordAttemptsVar).(int)),
		})
		if helper.IgnorePreconditionError(err) != nil {
			return diag.Errorf("failed to update default lockout policy: %v", err)
		}
		if resp != nil {
			id = resp.GetDetails().GetResourceOwner()
		}
	}
	if id == "" {
		resp, err := client.GetLockoutPolicy(ctx, &admin.GetLockoutPolicyRequest{})
		if err != nil {
			return diag.Errorf("failed to update default lockout policy: %v", err)
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

	client, err := helper.GetAdminClient(clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := client.GetLockoutPolicy(ctx, &admin.GetLockoutPolicyRequest{})
	if err != nil && helper.IgnoreIfNotFoundError(err) == nil {
		d.SetId("")
		return nil
	}
	if err != nil {
		return diag.Errorf("failed to get default lockout policy")
	}

	policy := resp.Policy
	set := map[string]interface{}{
		maxPasswordAttemptsVar: policy.GetMaxPasswordAttempts(),
	}

	for k, v := range set {
		if err := d.Set(k, v); err != nil {
			return diag.Errorf("failed to set %s of default lockout policy: %v", k, err)
		}
	}
	d.SetId(policy.GetDetails().GetResourceOwner())
	return nil
}
