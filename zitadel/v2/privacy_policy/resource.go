package privacy_policy

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func GetResource() *schema.Resource {
	return &schema.Resource{
		Description: "Resource representing the custom privacy policy of an organization.",
		Schema: map[string]*schema.Schema{
			orgIDVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Id for the organization",
				ForceNew:    true,
			},
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
		},
		CreateContext: create,
		DeleteContext: delete,
		ReadContext:   read,
		UpdateContext: update,
		Importer:      &schema.ResourceImporter{StateContext: schema.ImportStatePassthroughContext},
	}
}
