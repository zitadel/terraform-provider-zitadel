package lockout_policy

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
)

func GetResource() *schema.Resource {
	return &schema.Resource{
		Description: "Resource representing the custom lockout policy of an organization.",
		Schema: map[string]*schema.Schema{
			helper.OrgIDVar: helper.OrgIDResourceField,
			maxPasswordAttemptsVar: {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Maximum password check attempts before the account gets locked. Attempts are reset as soon as the password is entered correct or the password is reset.",
			},
			maxOTPAttemptsVar: {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Maximum OTP check attempts before the account gets locked. Attempts are reset as soon as the OTP is entered correct.",
			},
		},
		DeleteContext: delete,
		CreateContext: create,
		UpdateContext: update,
		ReadContext:   read,
		Importer:      helper.ImportWithOptionalOrg(),
	}
}
