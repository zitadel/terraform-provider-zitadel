package default_privacy_policy

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/helper"
)

func GetResource() *schema.Resource {
	return &schema.Resource{
		Description: "Resource representing the default privacy policy.",
		Schema: map[string]*schema.Schema{
			tosLinkVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "",
			},
			privacyLinkVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "",
			},
			helpLinkVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "",
			},
			supportEmailVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "",
			},
		},
		CreateContext: update,
		DeleteContext: delete,
		ReadContext:   read,
		UpdateContext: update,
		Importer:      &schema.ResourceImporter{StateContext: helper.ImportWithEmptyID()},
	}
}
