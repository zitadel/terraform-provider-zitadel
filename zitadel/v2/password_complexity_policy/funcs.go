package password_complexity_policy

import (
	"context"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/management"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/helper"
)

func delete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started create")

	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	org := d.Get(orgIDVar).(string)
	client, err := helper.GetManagementClient(clientinfo, org)
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.ResetPasswordComplexityPolicyToDefault(ctx, &management.ResetPasswordComplexityPolicyToDefaultRequest{})
	if err != nil {
		return diag.Errorf("failed to reset password complexity policy: %v", err)
	}
	return nil
}

func update(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started update")

	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	org := d.Get(orgIDVar).(string)
	client, err := helper.GetManagementClient(clientinfo, org)
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.UpdateCustomPasswordComplexityPolicy(ctx, &management.UpdateCustomPasswordComplexityPolicyRequest{
		MinLength:    uint64(d.Get(minLengthVar).(int)),
		HasUppercase: d.Get(hasUppercaseVar).(bool),
		HasLowercase: d.Get(hasLowercaseVar).(bool),
		HasNumber:    d.Get(hasNumberVar).(bool),
		HasSymbol:    d.Get(hasSymbolVar).(bool),
	})
	if err != nil {
		return diag.Errorf("failed to update password complexity policy: %v", err)
	}
	return nil
}

func create(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started create")

	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	org := d.Get(orgIDVar).(string)
	client, err := helper.GetManagementClient(clientinfo, org)
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.AddCustomPasswordComplexityPolicy(ctx, &management.AddCustomPasswordComplexityPolicyRequest{
		MinLength:    uint64(d.Get(minLengthVar).(int)),
		HasUppercase: d.Get(hasUppercaseVar).(bool),
		HasLowercase: d.Get(hasLowercaseVar).(bool),
		HasNumber:    d.Get(hasNumberVar).(bool),
		HasSymbol:    d.Get(hasSymbolVar).(bool),
	})
	if err != nil {
		return diag.Errorf("failed to create password complexity policy: %v", err)
	}
	d.SetId(org)
	return nil
}

func read(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started read")

	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	org := d.Get(orgIDVar).(string)
	client, err := helper.GetManagementClient(clientinfo, org)
	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := client.GetPasswordComplexityPolicy(ctx, &management.GetPasswordComplexityPolicyRequest{})
	if err != nil && helper.IgnoreIfNotFoundError(err) == nil {
		d.SetId("")
		return nil
	}
	if err != nil {
		return diag.Errorf("failed to get password complexity policy")
	}

	policy := resp.Policy
	if policy.GetIsDefault() == true {
		d.SetId("")
		return nil
	}
	set := map[string]interface{}{
		orgIDVar:        policy.GetDetails().GetResourceOwner(),
		minLengthVar:    policy.GetMinLength(),
		hasUppercaseVar: policy.GetHasUppercase(),
		hasLowercaseVar: policy.GetHasLowercase(),
		hasNumberVar:    policy.GetHasNumber(),
		hasSymbolVar:    policy.GetHasSymbol(),
	}

	for k, v := range set {
		if err := d.Set(k, v); err != nil {
			return diag.Errorf("failed to set %s of password complexity policy: %v", k, err)
		}
	}
	d.SetId(policy.GetDetails().GetResourceOwner())
	return nil
}
