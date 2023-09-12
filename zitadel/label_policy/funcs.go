package label_policy

import (
	"context"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/management"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/helper"
)

func delete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started delete")

	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := helper.GetManagementClient(clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.ResetLabelPolicyToDefault(helper.CtxWithOrgID(ctx, d), &management.ResetLabelPolicyToDefaultRequest{})
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

	org := helper.GetID(d, helper.OrgIDVar)
	client, err := helper.GetManagementClient(clientinfo)
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
		resp, err := client.UpdateCustomLabelPolicy(helper.CtxWithOrgID(ctx, d), &management.UpdateCustomLabelPolicyRequest{
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

	if d.HasChanges(LogoHashVar, LogoPathVar) {
		if err := helper.OrgFormFilePost(clientinfo, logoURL, d.Get(LogoPathVar).(string), org); err != nil {
			return diag.Errorf("failed to upload logo: %v", err)
		}
	}
	if d.HasChanges(LogoDarkHashVar, LogoDarkPathVar) {
		if err := helper.OrgFormFilePost(clientinfo, logoDarkURL, d.Get(LogoDarkPathVar).(string), org); err != nil {
			return diag.Errorf("failed to upload logo dark: %v", err)
		}
	}
	if d.HasChanges(IconHashVar, IconPathVar) {
		if err := helper.OrgFormFilePost(clientinfo, iconURL, d.Get(IconPathVar).(string), org); err != nil {
			return diag.Errorf("failed to upload icon: %v", err)
		}
	}
	if d.HasChanges(IconDarkHashVar, IconDarkPathVar) {
		if err := helper.OrgFormFilePost(clientinfo, iconDarkURL, d.Get(IconDarkPathVar).(string), org); err != nil {
			return diag.Errorf("failed to upload icon dark: %v", err)
		}
	}
	if d.HasChanges(FontHashVar, FontPathVar) {
		if err := helper.OrgFormFilePost(clientinfo, fontURL, d.Get(FontPathVar).(string), org); err != nil {
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
		LogoHashVar,
		LogoDarkHashVar,
		IconHashVar,
		IconDarkHashVar,
		FontHashVar,
	) {
		if d.Get(SetActiveVar).(bool) {
			if _, err := client.ActivateCustomLabelPolicy(helper.CtxWithOrgID(ctx, d), &management.ActivateCustomLabelPolicyRequest{}); err != nil {
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

	org := d.Get(helper.OrgIDVar).(string)
	client, err := helper.GetManagementClient(clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.AddCustomLabelPolicy(helper.CtxWithOrgID(ctx, d), &management.AddCustomLabelPolicyRequest{
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

	if d.Get(LogoHashVar) != "" && d.Get(LogoPathVar) != "" {
		if err := helper.OrgFormFilePost(clientinfo, logoURL, d.Get(LogoPathVar).(string), org); err != nil {
			return diag.Errorf("failed to upload logo: %v", err)
		}
	}
	if d.Get(LogoDarkHashVar) != "" && d.Get(LogoDarkPathVar) != "" {
		if err := helper.OrgFormFilePost(clientinfo, logoDarkURL, d.Get(LogoDarkPathVar).(string), org); err != nil {
			return diag.Errorf("failed to upload logo dark: %v", err)
		}
	}
	if d.Get(IconHashVar) != "" && d.Get(IconPathVar) != "" {
		if err := helper.OrgFormFilePost(clientinfo, iconURL, d.Get(IconPathVar).(string), org); err != nil {
			return diag.Errorf("failed to upload icon: %v", err)
		}
	}
	if d.Get(IconDarkHashVar) != "" && d.Get(IconDarkPathVar) != "" {
		if err := helper.OrgFormFilePost(clientinfo, iconDarkURL, d.Get(IconDarkPathVar).(string), org); err != nil {
			return diag.Errorf("failed to upload icon dark: %v", err)
		}
	}
	if d.Get(FontHashVar) != "" && d.Get(FontPathVar) != "" {
		if err := helper.OrgFormFilePost(clientinfo, fontURL, d.Get(FontPathVar).(string), org); err != nil {
			return diag.Errorf("failed to upload font: %v", err)
		}
	}

	if d.Get(SetActiveVar).(bool) {
		if _, err := client.ActivateCustomLabelPolicy(helper.CtxWithOrgID(ctx, d), &management.ActivateCustomLabelPolicyRequest{}); err != nil {
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

	client, err := helper.GetManagementClient(clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := client.GetPreviewLabelPolicy(helper.CtxWithOrgID(ctx, d), &management.GetPreviewLabelPolicyRequest{})
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
		helper.OrgIDVar:        policy.GetDetails().GetResourceOwner(),
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
