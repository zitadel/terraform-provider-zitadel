package privacy_policy

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/helper"
)

func GetResource() *schema.Resource {
	return &schema.Resource{
		Description: "Resource representing the custom privacy policy of an organization.",
		Schema: map[string]*schema.Schema{
			helper.OrgIDVar: helper.OrgIDResourceField,
			tosLinkVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "",
			},
			privacyLinkVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "",
			},
			HelpLinkVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "",
			},
			supportEmailVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "",
			},
		},
		CreateContext: create,
		DeleteContext: delete,
		ReadContext:   read,
		UpdateContext: update,
		Importer:      helper.ImportWithOptionalOrg(),
	}
}
