package notification_policy

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/helper"
)

func GetResource() *schema.Resource {
	return &schema.Resource{
		Description: "Resource representing the custom notification policy of an organization.",
		Schema: map[string]*schema.Schema{
			helper.OrgIDVar: helper.OrgIDResourceField,
			passwordChangeVar: {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Send notification if a user changes his password",
			},
		},
		ReadContext:   read,
		CreateContext: create,
		DeleteContext: delete,
		UpdateContext: update,
		Importer:      helper.ImportWithOptionalOrg(),
	}
}
