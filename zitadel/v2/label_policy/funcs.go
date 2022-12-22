package label_policy

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

	_, err = client.ResetLabelPolicyToDefault(ctx, &management.ResetLabelPolicyToDefaultRequest{})
	if err != nil {
		return diag.Errorf("failed to reset label policy: %v", err)
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
		resp, err := client.UpdateCustomLabelPolicy(ctx, &management.UpdateCustomLabelPolicyRequest{
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
		if err != nil {
			return diag.Errorf("failed to update label policy: %v", err)
		}
		d.SetId(resp.Details.ResourceOwner)
	}

	if d.HasChanges(logoHashVar, logoPathVar) {
		if err := helper.OrgFormFilePost(clientinfo, logoURL, d.Get(logoPathVar).(string), org); err != nil {
			return diag.Errorf("failed to upload logo: %v", err)
		}
	}
	if d.HasChanges(logoDarkHashVar, logoDarkPathVar) {
		if err := helper.OrgFormFilePost(clientinfo, logoDarkURL, d.Get(logoDarkPathVar).(string), org); err != nil {
			return diag.Errorf("failed to upload logo dark: %v", err)
		}
	}
	if d.HasChanges(iconHashVar, iconPathVar) {
		if err := helper.OrgFormFilePost(clientinfo, iconURL, d.Get(iconPathVar).(string), org); err != nil {
			return diag.Errorf("failed to upload icon: %v", err)
		}
	}
	if d.HasChanges(iconDarkHashVar, iconDarkPathVar) {
		if err := helper.OrgFormFilePost(clientinfo, iconDarkURL, d.Get(iconDarkPathVar).(string), org); err != nil {
			return diag.Errorf("failed to upload icon dark: %v", err)
		}
	}
	if d.HasChanges(fontHashVar, fontPathVar) {
		if err := helper.OrgFormFilePost(clientinfo, fontURL, d.Get(fontPathVar).(string), org); err != nil {
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
			if _, err := client.ActivateCustomLabelPolicy(ctx, &management.ActivateCustomLabelPolicyRequest{}); err != nil {
				return diag.Errorf("failed to activate label policy: %v", err)
			}
		}
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

	_, err = client.AddCustomLabelPolicy(ctx, &management.AddCustomLabelPolicyRequest{
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
	if err != nil {
		return diag.Errorf("failed to create label policy: %v", err)
	}
	d.SetId(org)

	if d.Get(logoHashVar) != "" && d.Get(logoPathVar) != "" {
		if err := helper.OrgFormFilePost(clientinfo, logoURL, d.Get(logoPathVar).(string), org); err != nil {
			return diag.Errorf("failed to upload logo: %v", err)
		}
	}
	if d.Get(logoDarkHashVar) != "" && d.Get(logoDarkPathVar) != "" {
		if err := helper.OrgFormFilePost(clientinfo, logoDarkURL, d.Get(logoDarkPathVar).(string), org); err != nil {
			return diag.Errorf("failed to upload logo dark: %v", err)
		}
	}
	if d.Get(iconHashVar) != "" && d.Get(iconPathVar) != "" {
		if err := helper.OrgFormFilePost(clientinfo, iconURL, d.Get(iconPathVar).(string), org); err != nil {
			return diag.Errorf("failed to upload icon: %v", err)
		}
	}
	if d.Get(iconDarkHashVar) != "" && d.Get(iconDarkPathVar) != "" {
		if err := helper.OrgFormFilePost(clientinfo, iconDarkURL, d.Get(iconDarkPathVar).(string), org); err != nil {
			return diag.Errorf("failed to upload icon dark: %v", err)
		}
	}
	if d.Get(fontHashVar) != "" && d.Get(fontPathVar) != "" {
		if err := helper.OrgFormFilePost(clientinfo, fontURL, d.Get(fontPathVar).(string), org); err != nil {
			return diag.Errorf("failed to upload font: %v", err)
		}
	}

	if d.Get(setActiveVar).(bool) {
		if _, err := client.ActivateCustomLabelPolicy(ctx, &management.ActivateCustomLabelPolicyRequest{}); err != nil {
			return diag.Errorf("failed to activate label policy: %v", err)
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

	org := d.Get(orgIDVar).(string)
	client, err := helper.GetManagementClient(clientinfo, org)
	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := client.GetPreviewLabelPolicy(ctx, &management.GetPreviewLabelPolicyRequest{})
	if err != nil && helper.IgnoreIfNotFoundError(err) == nil {
		d.SetId("")
		return nil
	}
	if err != nil {
		return diag.Errorf("failed to get label policy")
	}

	policy := resp.Policy
	if policy.GetIsDefault() == true {
		d.SetId("")
		return nil
	}
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
			return diag.Errorf("failed to set %s of label policy: %v", k, err)
		}
	}
	d.SetId(policy.GetDetails().GetResourceOwner())
	return nil
}
