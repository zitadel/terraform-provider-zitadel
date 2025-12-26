package default_lockout_policy

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
)

func GetResource() *schema.Resource {
	return &schema.Resource{
		Description: "Resource representing the default lockout policy.",
		Schema: map[string]*schema.Schema{
			MaxPasswordAttemptsVar: {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Maximum password check attempts before the account gets locked. Attempts are reset as soon as the password is entered correctly or the password is reset.",
			},
			MaxOTPAttemptsVar: {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Maximum OTP check attempts before the account gets locked. Attempts are reset as soon as the OTP is entered correctly.",
			},
		},
		DeleteContext: delete,
		CreateContext: update,
		UpdateContext: update,
		ReadContext:   read,
		Importer:      helper.ImportWithEmptyID(),
	}
}
