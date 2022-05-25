package v1

import (
	"context"
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

func GetLabelPolicyDatasource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			labelPolicyOrgIdVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Id for the organization",
			},
			labelPolicyPrimaryColor: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "hex value for primary color",
			},
			labelPolicyIsDefault: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "defines if the organisation's admin changed the policy",
			},
			labelPolicyHideLoginNameSuffix: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "hides the org suffix on the login form if the scope \"urn:zitadel:iam:org:domain:primary:{domainname}\" is set. Details about this scope in https://docs.zitadel.ch/concepts#Reserved_Scopes",
			},
			labelPolicyWarnColor: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "hex value for warn color",
			},
			labelPolicyBackgroundColor: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "hex value for background color",
			},
			labelPolicyFontColor: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "hex value for font color",
			},
			labelPolicyPrimaryColorDark: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "hex value for primary color dark theme",
			},
			labelPolicyBackgroundColorDark: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "hex value for background color dark theme",
			},
			labelPolicyWarnColorDark: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "hex value for warn color dark theme",
			},
			labelPolicyFontColorDark: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "hex value for font color dark theme",
			},
			labelPolicyDisableWatermark: {
				Type:        schema.TypeBool,
				Computed:    true,
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
	}
}

func readLabelPolicyOfOrg(ctx context.Context, policies *schema.Set, m interface{}, clientinfo *ClientInfo, org string) diag.Diagnostics {
	client, err := getManagementClient(clientinfo, org)
	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := client.GetLabelPolicy(ctx, &management2.GetLabelPolicyRequest{})
	if err != nil {
		return diag.Errorf("failed to get list of domains: %v", err)
	}

	policy := resp.Policy
	values := map[string]interface{}{
		labelPolicyOrgIdVar:            policy.GetDetails().GetResourceOwner(),
		labelPolicyPrimaryColor:        policy.GetPrimaryColor(),
		labelPolicyIsDefault:           policy.GetIsDefault(),
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
	policies.Add(values)
	return nil
}
