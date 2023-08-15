package password_complexity_policy

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/helper"
)

func GetResource() *schema.Resource {
	return &schema.Resource{
		Description: "Resource representing the custom password complexity policy of an organization.",
		Schema: map[string]*schema.Schema{
			helper.OrgIDVar: helper.OrgIDResourceField,
			minLengthVar: {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Minimal length for the password",
			},
			hasUppercaseVar: {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "defines if the password MUST contain an upper case letter",
			},
			hasLowercaseVar: {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "defines if the password MUST contain a lower case letter",
			},
			hasNumberVar: {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "defines if the password MUST contain a number",
			},
			hasSymbolVar: {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "defines if the password MUST contain a symbol. E.g. \"$\"",
			},
		},
		DeleteContext: delete,
		ReadContext:   read,
		CreateContext: create,
		UpdateContext: update,
		Importer:      &schema.ResourceImporter{StateContext: helper.ImportWithAttributesV5(helper.ImportOptionalOrgAttribute)},
	}
}
