package default_privacy_policy

import (
	"context"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/admin"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/helper"
)

func delete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "default privacy policy cannot be deleted")
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
	if d.HasChanges(tosLinkVar, privacyLinkVar, helpLinkVar) {
		resp, err := client.UpdatePrivacyPolicy(ctx, &admin.UpdatePrivacyPolicyRequest{
			TosLink:     d.Get(tosLinkVar).(string),
			PrivacyLink: d.Get(privacyLinkVar).(string),
			HelpLink:    d.Get(helpLinkVar).(string),
		})
		if helper.IgnorePreconditionError(err) != nil {
			return diag.Errorf("failed to update default privacy policy: %v", err)
		}
		if resp != nil {
			id = resp.GetDetails().GetResourceOwner()
		}
	}
	if id == "" {
		resp, err := client.GetPrivacyPolicy(ctx, &admin.GetPrivacyPolicyRequest{})
		if err != nil {
			return diag.Errorf("failed to update default privacy policy: %v", err)
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

	resp, err := client.GetPrivacyPolicy(ctx, &admin.GetPrivacyPolicyRequest{})
	if err != nil && helper.IgnoreIfNotFoundError(err) == nil {
		d.SetId("")
		return nil
	}
	if err != nil {
		return diag.Errorf("failed to default privacy policy")
	}

	policy := resp.Policy
	set := map[string]interface{}{
		tosLinkVar:     policy.GetTosLink(),
		privacyLinkVar: policy.GetPrivacyLink(),
		helpLinkVar:    policy.GetHelpLink(),
	}

	for k, v := range set {
		if err := d.Set(k, v); err != nil {
			return diag.Errorf("failed to set %s of default privacy policy: %v", k, err)
		}
	}
	d.SetId(policy.GetDetails().GetResourceOwner())
	return nil
}
