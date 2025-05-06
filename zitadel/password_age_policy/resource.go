package password_age_policy

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
)

func GetResource() *schema.Resource {
	return &schema.Resource{
		Description: "Resource representing the custom password age policy of an organization.",
		Schema: map[string]*schema.Schema{
			helper.OrgIDVar: helper.OrgIDResourceField,
			maxAgeDays: {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "amount of days after which a password will expire",
			},
			expireWarnDays: {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "amount of days after which the user should be notified of the upcoming expiry",
			},
		},
		DeleteContext: delete,
		ReadContext:   read,
		CreateContext: create,
		UpdateContext: update,
		Importer:      helper.ImportWithOptionalOrg(),
	}
}
