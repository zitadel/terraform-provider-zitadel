package domain_policy

import (
	"context"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/admin"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/management"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
)

func delete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started create")

	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := helper.GetAdminClient(ctx, clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}
	org := d.Get(helper.OrgIDVar).(string)

	_, err = client.ResetCustomDomainPolicyToDefault(helper.CtxWithOrgID(ctx, d), &admin.ResetCustomDomainPolicyToDefaultRequest{
		OrgId: org,
	})
	if err != nil {
		return diag.Errorf("failed to reset domain policy: %v", err)
	}
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
	org := helper.GetID(d, helper.OrgIDVar)
	_, err = client.UpdateCustomDomainPolicy(ctx, &admin.UpdateCustomDomainPolicyRequest{
		OrgId:                                  org,
		UserLoginMustBeDomain:                  d.Get(UserLoginMustBeDomainVar).(bool),
		ValidateOrgDomains:                     d.Get(validateOrgDomainVar).(bool),
		SmtpSenderAddressMatchesInstanceDomain: d.Get(smtpSenderVar).(bool),
	})
	if err != nil {
		return diag.Errorf("failed to update domain policy: %v", err)
	}
	d.SetId(org)
	return nil
}

func create(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started create")

	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := helper.GetAdminClient(ctx, clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}
	org := helper.GetID(d, helper.OrgIDVar)
	_, err = client.AddCustomDomainPolicy(ctx, &admin.AddCustomDomainPolicyRequest{
		OrgId:                                  org,
		UserLoginMustBeDomain:                  d.Get(UserLoginMustBeDomainVar).(bool),
		ValidateOrgDomains:                     d.Get(validateOrgDomainVar).(bool),
		SmtpSenderAddressMatchesInstanceDomain: d.Get(smtpSenderVar).(bool),
	})
	if err != nil {
		return diag.Errorf("failed to create domain policy: %v", err)
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

	client, err := helper.GetManagementClient(ctx, clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}

	org := helper.GetID(d, helper.OrgIDVar)
	resp, err := client.GetDomainPolicy(helper.CtxSetOrgID(ctx, org), &management.GetDomainPolicyRequest{})
	if err != nil && helper.IgnoreIfNotFoundError(err) == nil {
		d.SetId("")
		return nil
	}
	if err != nil {
		return diag.Errorf("failed to get domain policy")
	}

	policy := resp.Policy
	if policy.GetIsDefault() == true {
		d.SetId("")
		return nil
	}
	set := map[string]interface{}{
		helper.OrgIDVar:          policy.GetDetails().GetResourceOwner(),
		UserLoginMustBeDomainVar: policy.GetUserLoginMustBeDomain(),
		validateOrgDomainVar:     policy.GetValidateOrgDomains(),
		smtpSenderVar:            policy.GetSmtpSenderAddressMatchesInstanceDomain(),
	}

	for k, v := range set {
		if err := d.Set(k, v); err != nil {
			return diag.Errorf("failed to set %s of domain: %v", k, err)
		}
	}
	d.SetId(policy.GetDetails().GetResourceOwner())
	return nil
}
