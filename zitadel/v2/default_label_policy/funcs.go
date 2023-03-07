package default_label_policy

import (
	"context"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/admin"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/helper"
)

func delete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "default label policy cannot be deleted")
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
	if d.HasChanges(
		primaryColorVar,
		hideLoginNameSuffixVar,
		warnColorVar,
		backgroundColorVar,
		fontColorVar,
		primaryColorDarkVar,
		backgroundColorDarkVar,
		warnColorDarkVar,
		fontColorDarkVar,
		disableWatermarkVar,
	) {
		resp, err := client.UpdateLabelPolicy(ctx, &admin.UpdateLabelPolicyRequest{
			PrimaryColor:        d.Get(primaryColorVar).(string),
			HideLoginNameSuffix: d.Get(hideLoginNameSuffixVar).(bool),
			WarnColor:           d.Get(warnColorVar).(string),
			BackgroundColor:     d.Get(backgroundColorVar).(string),
			FontColor:           d.Get(fontColorVar).(string),
			PrimaryColorDark:    d.Get(primaryColorDarkVar).(string),
			BackgroundColorDark: d.Get(backgroundColorDarkVar).(string),
			WarnColorDark:       d.Get(warnColorDarkVar).(string),
			FontColorDark:       d.Get(fontColorDarkVar).(string),
			DisableWatermark:    d.Get(disableWatermarkVar).(bool),
		})
		if helper.IgnorePreconditionError(err) != nil {
			return diag.Errorf("failed to update default label policy: %v", err)
		}
		if resp != nil {
			id = resp.Details.ResourceOwner
		}
	}
	if id == "" {
		resp, err := client.GetLabelPolicy(ctx, &admin.GetLabelPolicyRequest{})
		if err != nil {
			return diag.Errorf("failed to update default label policy: %v", err)
		}
		id = resp.GetPolicy().GetDetails().GetResourceOwner()
	}
	d.SetId(id)

	if d.HasChanges(logoHashVar, logoPathVar) {
		if err := helper.InstanceFormFilePost(clientinfo, logoURL, d.Get(logoPathVar).(string)); err != nil {
			return diag.Errorf("failed to upload logo: %v", err)
		}
	}
	if d.HasChanges(logoDarkHashVar, logoDarkPathVar) {
		if err := helper.InstanceFormFilePost(clientinfo, logoDarkURL, d.Get(logoDarkPathVar).(string)); err != nil {
			return diag.Errorf("failed to upload logo dark: %v", err)
		}
	}
	if d.HasChanges(iconHashVar, iconPathVar) {
		if err := helper.InstanceFormFilePost(clientinfo, iconURL, d.Get(iconPathVar).(string)); err != nil {
			return diag.Errorf("failed to upload icon: %v", err)
		}
	}
	if d.HasChanges(iconDarkHashVar, iconDarkPathVar) {
		if err := helper.InstanceFormFilePost(clientinfo, iconDarkURL, d.Get(iconDarkPathVar).(string)); err != nil {
			return diag.Errorf("failed to upload icon dark: %v", err)
		}
	}
	if d.HasChanges(fontHashVar, fontPathVar) {
		if err := helper.InstanceFormFilePost(clientinfo, fontURL, d.Get(fontPathVar).(string)); err != nil {
			return diag.Errorf("failed to upload font: %v", err)
		}
	}

	if d.HasChanges(
		primaryColorVar,
		hideLoginNameSuffixVar,
		warnColorVar,
		backgroundColorVar,
		fontColorVar,
		primaryColorDarkVar,
		backgroundColorDarkVar,
		warnColorDarkVar,
		fontColorDarkVar,
		disableWatermarkVar,
		logoHashVar,
		logoDarkHashVar,
		iconHashVar,
		iconDarkHashVar,
		fontHashVar,
	) {
		if d.Get(setActiveVar).(bool) {
			if _, err := client.ActivateLabelPolicy(ctx, &admin.ActivateLabelPolicyRequest{}); err != nil {
				return diag.Errorf("failed to activate default label policy: %v", err)
			}
		}
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

	resp, err := client.GetLabelPolicy(ctx, &admin.GetLabelPolicyRequest{})
	if err != nil && helper.IgnoreIfNotFoundError(err) == nil {
		d.SetId("")
		return nil
	}
	if err != nil {
		return diag.Errorf("failed to get default label policy")
	}

	policy := resp.Policy
	set := map[string]interface{}{
		primaryColorVar:        policy.GetPrimaryColor(),
		hideLoginNameSuffixVar: policy.GetHideLoginNameSuffix(),
		warnColorVar:           policy.GetWarnColor(),
		backgroundColorVar:     policy.GetBackgroundColor(),
		fontColorVar:           policy.GetFontColor(),
		primaryColorDarkVar:    policy.GetPrimaryColorDark(),
		backgroundColorDarkVar: policy.GetBackgroundColorDark(),
		warnColorDarkVar:       policy.GetWarnColorDark(),
		fontColorDarkVar:       policy.GetFontColorDark(),
		disableWatermarkVar:    policy.GetDisableWatermark(),
		logoURLVar:             policy.GetLogoUrl(),
		iconURLVar:             policy.GetIconUrl(),
		logoURLDarkVar:         policy.GetLogoUrlDark(),
		iconURLDarkVar:         policy.GetIconUrlDark(),
		fontURLVar:             policy.GetFontUrl(),
	}

	for k, v := range set {
		if err := d.Set(k, v); err != nil {
			return diag.Errorf("failed to set %s of default label policy: %v", k, err)
		}
	}
	d.SetId(policy.GetDetails().GetResourceOwner())
	return nil
}
