package default_password_complexity_policy

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func GetResource() *schema.Resource {
	return &schema.Resource{
		Description: "Resource representing the default password complexity policy.",
		Schema: map[string]*schema.Schema{
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
		CreateContext: update,
		UpdateContext: update,
		Importer:      &schema.ResourceImporter{StateContext: schema.ImportStatePassthroughContext},
	}
}
