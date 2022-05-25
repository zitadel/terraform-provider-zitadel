package v2

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

const (
	privacyPolicyOrgIdVar    = "org_id"
	privacyPolicyTOSLink     = "tos_link"
	privacyPolicyPrivacyLink = "privacy_link"
	privacyPolicyIsDefault   = "is_default"
	privacyPolicyHelpLink    = "help_link"
)

func GetPrivacyPolicy() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			privacyPolicyOrgIdVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Id for the organization",
			},
			privacyPolicyTOSLink: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "",
			},
			privacyPolicyPrivacyLink: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "",
			},
			privacyPolicyIsDefault: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "",
			},
			privacyPolicyHelpLink: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "",
			},
		},
	}
}
