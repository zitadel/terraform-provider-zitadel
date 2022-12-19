package default_domain_policy

import (
	"context"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/admin"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/helper"
)

func delete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "default domain policy cannot be deleted")
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

	resp, err := client.UpdateDomainPolicy(ctx, &admin.UpdateDomainPolicyRequest{
		UserLoginMustBeDomain:                  d.Get(userLoginMustBeDomainVar).(bool),
		ValidateOrgDomains:                     d.Get(validateOrgDomainVar).(bool),
		SmtpSenderAddressMatchesInstanceDomain: d.Get(smtpSenderVar).(bool),
	})
	if helper.IgnorePreconditionError(err) != nil {
		return diag.Errorf("failed to update default domain policy: %v", err)
	}
	if resp != nil {
		d.SetId(resp.GetDetails().GetResourceOwner())
	} else {
		resp, err := client.GetDomainPolicy(ctx, &admin.GetDomainPolicyRequest{})
		if err != nil {
			return diag.Errorf("failed to update default domain policy: %v", err)
		}
		d.SetId(resp.GetPolicy().GetDetails().GetResourceOwner())
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

	resp, err := client.GetDomainPolicy(ctx, &admin.GetDomainPolicyRequest{})
	if err != nil && helper.IgnoreIfNotFoundError(err) == nil {
		d.SetId("")
		return nil
	}
	if err != nil {
		return diag.Errorf("failed to get default domain policy")
	}
	policy := resp.Policy
	set := map[string]interface{}{
		userLoginMustBeDomainVar: policy.GetUserLoginMustBeDomain(),
		validateOrgDomainVar:     policy.GetValidateOrgDomains(),
	}

	for k, v := range set {
		if err := d.Set(k, v); err != nil {
			return diag.Errorf("failed to set %s of default domain policy: %v", k, err)
		}
	}
	d.SetId(policy.GetDetails().GetResourceOwner())
	return nil
}
