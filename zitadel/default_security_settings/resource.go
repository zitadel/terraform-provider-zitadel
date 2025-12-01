package default_security_settings

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
)

func GetResource() *schema.Resource {
	return &schema.Resource{
		Description: "Resource representing the default security settings.",
		Schema: map[string]*schema.Schema{
			EnableImpersonationVar: {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Enable impersonation for the instance",
			},
			embeddedIframeEnabledVar: {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Enable embedding in iframes",
			},
			embeddedIframeAllowedOriginsVar: {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				Description: "Origins allowed to embed ZITADEL in an iframe",
			},
		},
		ReadContext:   read,
		CreateContext: update,
		DeleteContext: delete,
		UpdateContext: update,
		Importer:      helper.ImportWithEmptyID(),
	}
}
