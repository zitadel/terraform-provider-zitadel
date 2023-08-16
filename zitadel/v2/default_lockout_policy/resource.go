package default_lockout_policy

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/helper"
)

func GetResource() *schema.Resource {
	return &schema.Resource{
		Description: "Resource representing the default lockout policy.",
		Schema: map[string]*schema.Schema{
			maxPasswordAttemptsVar: {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Maximum password check attempts before the account gets locked. Attempts are reset as soon as the password is entered correctly or the password is reset.",
			},
		},
		DeleteContext: delete,
		CreateContext: update,
		UpdateContext: update,
		ReadContext:   read,
		Importer:      &schema.ResourceImporter{StateContext: helper.ImportWithEmptyID()},
	}
}
