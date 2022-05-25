package v2

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

const (
	passwordCompPolicyOrgIdVar     = "org_id"
	passwordCompPolicyMinLength    = "min_length"
	passwordCompPolicyHasUppercase = "has_uppercase"
	passwordCompPolicyHasLowercase = "has_lowercase"
	passwordCompPolicyHasNumber    = "has_number"
	passwordCompPolicyHasSymbol    = "has_symbol"
	passwordCompPolicyIsDefault    = "is_default"
)

func GetPasswordComplexityPolicy() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			passwordCompPolicyOrgIdVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Id for the organization",
			},
			passwordCompPolicyMinLength: {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Minimal length for the password",
			},
			passwordCompPolicyHasUppercase: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "defines if the password MUST contain an upper case letter",
			},
			passwordCompPolicyHasLowercase: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "defines if the password MUST contain a lower case letter",
			},
			passwordCompPolicyHasNumber: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "defines if the password MUST contain a number",
			},
			passwordCompPolicyHasSymbol: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "defines if the password MUST contain a symbol. E.g. \"$\"",
			},
			passwordCompPolicyIsDefault: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "defines if the organisation's admin changed the policy",
			},
		},
	}
}
