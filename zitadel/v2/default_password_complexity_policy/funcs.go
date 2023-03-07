package default_password_complexity_policy

import (
	"context"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/admin"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/helper"
)

func delete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "default password complexity policy cannot be deleted")
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
	if d.HasChanges(minLengthVar, hasUppercaseVar, hasLowercaseVar, hasNumberVar, hasSymbolVar) {
		resp, err := client.UpdatePasswordComplexityPolicy(ctx, &admin.UpdatePasswordComplexityPolicyRequest{
			MinLength:    uint32(d.Get(minLengthVar).(int)),
			HasUppercase: d.Get(hasUppercaseVar).(bool),
			HasLowercase: d.Get(hasLowercaseVar).(bool),
			HasNumber:    d.Get(hasNumberVar).(bool),
			HasSymbol:    d.Get(hasSymbolVar).(bool),
		})
		if helper.IgnorePreconditionError(err) != nil {
			return diag.Errorf("failed to update default password complexity policy: %v", err)
		}
		if resp != nil {
			id = resp.GetDetails().GetResourceOwner()
		}
	}
	if id == "" {
		resp, err := client.GetPasswordComplexityPolicy(ctx, &admin.GetPasswordComplexityPolicyRequest{})
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

	client, err := helper.GetAdminClient(clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := client.GetPasswordComplexityPolicy(ctx, &admin.GetPasswordComplexityPolicyRequest{})
	if err != nil && helper.IgnoreIfNotFoundError(err) == nil {
		d.SetId("")
		return nil
	}
	if err != nil {
		return diag.Errorf("failed to get default password complexity policy")
	}

	policy := resp.Policy
	set := map[string]interface{}{
		minLengthVar:    policy.GetMinLength(),
		hasUppercaseVar: policy.GetHasUppercase(),
		hasLowercaseVar: policy.GetHasLowercase(),
		hasNumberVar:    policy.GetHasNumber(),
		hasSymbolVar:    policy.GetHasSymbol(),
	}

	for k, v := range set {
		if err := d.Set(k, v); err != nil {
			return diag.Errorf("failed to set %s of default password complexity policy: %v", k, err)
		}
	}
	d.SetId(policy.GetDetails().GetResourceOwner())
	return nil
}
