package password_age_policy

import (
	"context"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/management"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
)

func delete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started create")

	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := helper.GetManagementClient(ctx, clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.ResetPasswordAgePolicyToDefault(helper.CtxWithID(ctx, d), &management.ResetPasswordAgePolicyToDefaultRequest{})
	if err != nil {
		return diag.Errorf("failed to reset password age policy: %v", err)
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

	req := management.UpdateCustomPasswordAgePolicyRequest{
		MaxAgeDays:     uint32(d.Get(maxAgeDays).(int)),
		ExpireWarnDays: uint32(d.Get(expireWarnDays).(int)),
	}

	_, err = client.UpdateCustomPasswordAgePolicy(helper.CtxWithID(ctx, d), &req)
	if err != nil {
		return diag.Errorf("failed to update password age policy: %v", err)
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

	_, err = client.AddCustomPasswordAgePolicy(helper.CtxWithID(ctx, d), &management.AddCustomPasswordAgePolicyRequest{
		MaxAgeDays:     uint32(d.Get(maxAgeDays).(int)),
		ExpireWarnDays: uint32(d.Get(expireWarnDays).(int)),
	})

	if err != nil {
		return diag.Errorf("failed to create password age policy: %v", err)
	}

	org := d.Get(helper.OrgIDVar).(string)
	d.SetId(org)

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

	resp, err := client.GetPasswordAgePolicy(helper.CtxWithID(ctx, d), &management.GetPasswordAgePolicyRequest{})
	if err != nil && helper.IgnoreIfNotFoundError(err) == nil {
		d.SetId("")
		return nil
	}

	if err != nil {
		return diag.Errorf("failed to get password age policy: %v", err)
	}

	policy := resp.Policy
	if policy.GetIsDefault() == true {
		d.SetId("")
		return nil
	}

	set := map[string]interface{}{
		helper.OrgIDVar: policy.GetDetails().GetResourceOwner(),
		maxAgeDays:      policy.GetMaxAgeDays(),
		expireWarnDays:  policy.GetExpireWarnDays(),
	}

	for k, v := range set {
		if err := d.Set(k, v); err != nil {
			return diag.Errorf("failed to set %s of password age policy: %v", k, err)
		}
	}

	d.SetId(policy.GetDetails().GetResourceOwner())
	return nil
}
