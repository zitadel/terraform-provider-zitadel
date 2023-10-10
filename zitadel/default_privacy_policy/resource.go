package default_privacy_policy

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/helper"
)

func GetResource() *schema.Resource {
	return &schema.Resource{
		Description: "Resource representing the default privacy policy.",
		Schema: map[string]*schema.Schema{
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
		CreateContext: update,
		DeleteContext: delete,
		ReadContext:   read,
		UpdateContext: update,
		Importer:      helper.ImportWithEmptyID(),
	}
}
