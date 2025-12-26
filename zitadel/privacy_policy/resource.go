package privacy_policy

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
)

func GetResource() *schema.Resource {
	return &schema.Resource{
		Description: "Resource representing the custom privacy policy of an organization.",
		Schema: map[string]*schema.Schema{
			helper.OrgIDVar: helper.OrgIDResourceField,
			tosLinkVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Link to the Terms of Service.",
			},
			privacyLinkVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Link to the Privacy Policy.",
			},
			HelpLinkVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Link to the Help/Manual page.",
			},
			supportEmailVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Help / support email address.",
			},
			DocsLinkVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Link to documentation to be shown in the console.",
			},
			CustomLinkVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Link to an external resource that will be available to users in the console.",
			},
			CustomLinkTextVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The button text that would be shown in console pointing to custom link.",
			},
		},
		CreateContext: create,
		DeleteContext: delete,
		ReadContext:   read,
		UpdateContext: update,
		Importer:      helper.ImportWithOptionalOrg(),
	}
}
