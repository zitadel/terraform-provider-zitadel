package default_notification_policy

import (
	"context"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/admin"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/helper"
)

func delete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "default notification policy cannot be deleted")
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

	if d.HasChanges(passwordChangeVar) {
		resp, err := client.UpdateNotificationPolicy(ctx, &admin.UpdateNotificationPolicyRequest{
			PasswordChange: d.Get(passwordChangeVar).(bool),
		})
		if helper.IgnorePreconditionError(err) != nil {
			return diag.Errorf("failed to update default notification policy: %v", err)
		}
		if resp != nil {
			d.SetId(resp.GetDetails().GetResourceOwner())
			return nil
		}
	}

	resp, err := client.GetNotificationPolicy(ctx, &admin.GetNotificationPolicyRequest{})
	if err != nil {
		return diag.Errorf("failed to update default notification policy: %v", err)
	}
	d.SetId(resp.GetPolicy().GetDetails().GetResourceOwner())
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

	resp, err := client.GetNotificationPolicy(ctx, &admin.GetNotificationPolicyRequest{})
	if err != nil && helper.IgnoreIfNotFoundError(err) == nil {
		d.SetId("")
		return nil
	}
	if err != nil {
		return diag.Errorf("failed to get default notification policy")
	}
	policy := resp.Policy
	set := map[string]interface{}{
		passwordChangeVar: policy.GetPasswordChange(),
	}

	for k, v := range set {
		if err := d.Set(k, v); err != nil {
			return diag.Errorf("failed to set %s of default notification policy: %v", k, err)
		}
	}
	d.SetId(policy.GetDetails().GetResourceOwner())
	return nil
}
