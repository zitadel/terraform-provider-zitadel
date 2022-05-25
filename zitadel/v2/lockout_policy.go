package v2

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

const (
	lockoutPolicyOrgIdVar            = "org_id"
	lockoutPolicyMaxPasswordAttempts = "user_login"
	lockoutPolicyIsDefault           = "is_default"
)

func GetLockoutPolicy() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			lockoutPolicyOrgIdVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Id for the organization",
			},
			lockoutPolicyMaxPasswordAttempts: {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Maximum password check attempts before the account gets locked. Attempts are reset as soon as the password is entered correct or the password is reset.",
			},
			lockoutPolicyIsDefault: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "defines if the organisation's admin changed the policy",
			},
		},
	}
}
