package default_security_settings

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
)

func GetResource() *schema.Resource {
	return &schema.Resource{
		Description: "Resource representing the default security settings.",
		Schema: map[string]*schema.Schema{
			enableImpersonationVar: {
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
			dynamicClientRegistrationEnabledVar: {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
				Description: "Serve the OAuth 2.0 Dynamic Client Registration endpoint (RFC 7591) and advertise it as " +
					"`registration_endpoint` in the discovery document. If not set, the value configured in ZITADEL is kept.",
			},
			dynamicClientRegistrationAllowUnauthenticatedVar: {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
				Description: "Allow clients to register dynamically without a bearer token, as the Model Context Protocol (MCP) " +
					"authorization flow requires. Such clients are created in the instance's default organization. Only has an " +
					"effect if `" + dynamicClientRegistrationEnabledVar + "` is true. If false, registration requires an access " +
					"token of a user with the `project.app.register_dynamic` permission. If not set, the value configured in " +
					"ZITADEL is kept.",
			},
		},
		ReadContext:   read,
		CreateContext: update,
		DeleteContext: delete,
		UpdateContext: update,
		Importer:      helper.ImportWithEmptyID(),
	}
}
