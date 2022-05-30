package v2

import (
	"context"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	management2 "github.com/zitadel/zitadel-go/pkg/client/zitadel/management"
)

const (
	labelPolicyOrgIdVar            = "org_id"
	labelPolicyPrimaryColor        = "primary_color"
	labelPolicyIsDefault           = "is_default"
	labelPolicyHideLoginNameSuffix = "hide_login_name_suffix"
	labelPolicyWarnColor           = "warn_color"
	labelPolicyBackgroundColor     = "background_color"
	labelPolicyFontColor           = "font_color"
	labelPolicyPrimaryColorDark    = "primary_color_dark"
	labelPolicyBackgroundColorDark = "background_color_dark"
	labelPolicyWarnColorDark       = "warn_color_dark"
	labelPolicyFontColorDark       = "font_color_dark"
	labelPolicyDisableWatermark    = "disable_watermark"
	labelPolicyLogoURL             = "logo_url"
	labelPolicyIconURL             = "icon_url"
	labelPolicyLogoURLDark         = "logo_url_dark"
	labelPolicyIconURLDark         = "icon_url_dark"
	labelPolicyFontURL             = "font_url"
)

func GetLabelPolicy() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			labelPolicyOrgIdVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Id for the organization",
				ForceNew:    true,
			},
			labelPolicyPrimaryColor: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "hex value for primary color",
			},
			labelPolicyIsDefault: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "defines if the organisation's admin changed the policy",
			},
			labelPolicyHideLoginNameSuffix: {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "hides the org suffix on the login form if the scope \"urn:zitadel:iam:org:domain:primary:{domainname}\" is set. Details about this scope in https://docs.zitadel.ch/concepts#Reserved_Scopes",
			},
			labelPolicyWarnColor: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "hex value for warn color",
			},
			labelPolicyBackgroundColor: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "hex value for background color",
			},
			labelPolicyFontColor: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "hex value for font color",
			},
			labelPolicyPrimaryColorDark: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "hex value for primary color dark theme",
			},
			labelPolicyBackgroundColorDark: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "hex value for background color dark theme",
			},
			labelPolicyWarnColorDark: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "hex value for warn color dark theme",
			},
			labelPolicyFontColorDark: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "hex value for font color dark theme",
			},
			labelPolicyDisableWatermark: {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "",
			},
			labelPolicyLogoURL: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "",
			},
			labelPolicyIconURL: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "",
			},
			labelPolicyLogoURLDark: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "",
			},
			labelPolicyIconURLDark: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "",
			},
			labelPolicyFontURL: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "",
			},
		},
		ReadContext:   readLabelPolicy,
		CreateContext: createLabelPolicy,
		DeleteContext: deleteLabelPolicy,
		UpdateContext: updateLabelPolicy,
	}
}

func deleteLabelPolicy(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started create")

	clientinfo, ok := m.(*ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	org := d.Get(labelPolicyOrgIdVar).(string)
	client, err := getManagementClient(clientinfo, org)
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.ResetLabelPolicyToDefault(ctx, &management2.ResetLabelPolicyToDefaultRequest{})
	if err != nil {
		return diag.Errorf("failed to reset label policy: %v", err)
	}
	d.SetId(org)
	return nil
}

func updateLabelPolicy(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started update")

	clientinfo, ok := m.(*ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	org := d.Get(labelPolicyOrgIdVar).(string)
	client, err := getManagementClient(clientinfo, org)
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.UpdateCustomLabelPolicy(ctx, &management2.UpdateCustomLabelPolicyRequest{
		PrimaryColor:        d.Get(labelPolicyPrimaryColor).(string),
		HideLoginNameSuffix: d.Get(labelPolicyHideLoginNameSuffix).(bool),
		WarnColor:           d.Get(labelPolicyWarnColor).(string),
		BackgroundColor:     d.Get(labelPolicyBackgroundColor).(string),
		FontColor:           d.Get(labelPolicyFontColor).(string),
		PrimaryColorDark:    d.Get(labelPolicyPrimaryColorDark).(string),
		BackgroundColorDark: d.Get(labelPolicyBackgroundColorDark).(string),
		WarnColorDark:       d.Get(labelPolicyWarnColorDark).(string),
		FontColorDark:       d.Get(labelPolicyFontColorDark).(string),
		DisableWatermark:    d.Get(labelPolicyDisableWatermark).(bool),
	})
	if err != nil {
		return diag.Errorf("failed to update label policy: %v", err)
	}
	d.SetId(org)
	return nil
}

func createLabelPolicy(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started create")

	clientinfo, ok := m.(*ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	org := d.Get(labelPolicyOrgIdVar).(string)
	client, err := getManagementClient(clientinfo, org)
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.AddCustomLabelPolicy(ctx, &management2.AddCustomLabelPolicyRequest{
		PrimaryColor:        d.Get(labelPolicyPrimaryColor).(string),
		HideLoginNameSuffix: d.Get(labelPolicyHideLoginNameSuffix).(bool),
		WarnColor:           d.Get(labelPolicyWarnColor).(string),
		BackgroundColor:     d.Get(labelPolicyBackgroundColor).(string),
		FontColor:           d.Get(labelPolicyFontColor).(string),
		PrimaryColorDark:    d.Get(labelPolicyPrimaryColorDark).(string),
		BackgroundColorDark: d.Get(labelPolicyBackgroundColorDark).(string),
		WarnColorDark:       d.Get(labelPolicyWarnColorDark).(string),
		FontColorDark:       d.Get(labelPolicyFontColorDark).(string),
		DisableWatermark:    d.Get(labelPolicyDisableWatermark).(bool),
	})
	if err != nil {
		return diag.Errorf("failed to create label policy: %v", err)
	}
	d.SetId(org)
	return nil
}

func readLabelPolicy(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started read")

	clientinfo, ok := m.(*ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	org := d.Get(domainPolicyOrgIdVar).(string)
	client, err := getManagementClient(clientinfo, org)
	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := client.GetLabelPolicy(ctx, &management2.GetLabelPolicyRequest{})
	if err != nil {
		return diag.Errorf("failed to get domain policy: %v", err)
	}

	policy := resp.Policy
	set := map[string]interface{}{
		labelPolicyOrgIdVar:            policy.GetDetails().GetResourceOwner(),
		labelPolicyIsDefault:           policy.GetIsDefault(),
		labelPolicyPrimaryColor:        policy.GetPrimaryColor(),
		labelPolicyHideLoginNameSuffix: policy.GetHideLoginNameSuffix(),
		labelPolicyWarnColor:           policy.GetWarnColor(),
		labelPolicyBackgroundColor:     policy.GetBackgroundColor(),
		labelPolicyFontColor:           policy.GetFontColor(),
		labelPolicyPrimaryColorDark:    policy.GetPrimaryColorDark(),
		labelPolicyBackgroundColorDark: policy.GetBackgroundColorDark(),
		labelPolicyWarnColorDark:       policy.GetWarnColorDark(),
		labelPolicyFontColorDark:       policy.GetFontColorDark(),
		labelPolicyDisableWatermark:    policy.GetDisableWatermark(),
		labelPolicyLogoURL:             policy.GetLogoUrl(),
		labelPolicyIconURL:             policy.GetIconUrl(),
		labelPolicyLogoURLDark:         policy.GetLogoUrlDark(),
		labelPolicyIconURLDark:         policy.GetIconUrlDark(),
		labelPolicyFontURL:             policy.GetFontUrl(),
	}

	for k, v := range set {
		if err := d.Set(k, v); err != nil {
			return diag.Errorf("failed to set %s of label policy: %v", k, err)
		}
	}
	d.SetId(policy.GetDetails().GetResourceOwner())
	return nil
}
