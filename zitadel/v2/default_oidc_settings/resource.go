package default_oidc_settings

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/helper"
)

func GetResource() *schema.Resource {
	return &schema.Resource{
		Description: "Resource representing the default oidc settings.",
		Schema: map[string]*schema.Schema{
			accessTokenLifetimeVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "lifetime duration of access tokens",
			},
			idTokenLifetimeVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "lifetime duration of id tokens",
			},
			refreshTokenExpirationVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "expiration duration of refresh tokens",
			},
			refreshTokenIdleExpirationVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "expiration duration of idle refresh tokens",
			},
		},
		CreateContext: update,
		UpdateContext: update,
		DeleteContext: delete,
		ReadContext:   read,
		Importer:      &schema.ResourceImporter{StateContext: helper.ImportWithEmptyID()},
	}
}
