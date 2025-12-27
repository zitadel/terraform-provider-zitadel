package notification_policy

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
	_, err = client.ResetNotificationPolicyToDefault(helper.CtxWithID(ctx, d), &management.ResetNotificationPolicyToDefaultRequest{})
	if err != nil {
		return diag.Errorf("failed to reset notification policy: %v", err)
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
	org := helper.GetID(d, helper.OrgIDVar)
	if d.HasChanges(passwordChangeVar) {
		_, err = client.UpdateCustomNotificationPolicy(helper.CtxWithID(ctx, d), &management.UpdateCustomNotificationPolicyRequest{
			PasswordChange: d.Get(passwordChangeVar).(bool),
		})
		if err != nil {
			return diag.Errorf("failed to update notification policy: %v", err)
		}
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
	client, err := helper.GetManagementClient(ctx, clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}
	org := d.Get(helper.OrgIDVar).(string)
	_, err = client.AddCustomNotificationPolicy(helper.CtxWithID(ctx, d), &management.AddCustomNotificationPolicyRequest{
		PasswordChange: d.Get(passwordChangeVar).(bool),
	})
	if err != nil {
		return diag.Errorf("failed to create notification policy: %v", err)
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
	resp, err := client.GetNotificationPolicy(helper.CtxWithID(ctx, d), &management.GetNotificationPolicyRequest{})
	if err != nil && helper.IgnoreIfNotFoundError(err) == nil {
		d.SetId("")
		return nil
	}
	if err != nil {
		return diag.Errorf("failed to get notification policy: %v", err)
	}
	policy := resp.Policy
	if policy.GetIsDefault() == true {
		d.SetId("")
		return nil
	}
	set := map[string]interface{}{
		helper.OrgIDVar:   policy.GetDetails().GetResourceOwner(),
		passwordChangeVar: policy.GetPasswordChange(),
	}
	for k, v := range set {
		if err := d.Set(k, v); err != nil {
			return diag.Errorf("failed to set %s of notification: %v", k, err)
		}
	}
	d.SetId(policy.GetDetails().GetResourceOwner())
	return nil
}
