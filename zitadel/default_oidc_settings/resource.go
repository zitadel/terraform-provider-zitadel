package default_oidc_settings

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
)

func GetResource() *schema.Resource {
	return &schema.Resource{
		Description: "Resource representing the default oidc settings.",
		Schema: map[string]*schema.Schema{
			accessTokenLifetimeVar: {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "lifetime duration of access tokens",
				DiffSuppressFunc: helper.DurationDiffSuppress,
			},
			idTokenLifetimeVar: {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "lifetime duration of id tokens",
				DiffSuppressFunc: helper.DurationDiffSuppress,
			},
			RefreshTokenExpirationVar: {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "expiration duration of refresh tokens",
				DiffSuppressFunc: helper.DurationDiffSuppress,
			},
			refreshTokenIdleExpirationVar: {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "expiration duration of idle refresh tokens",
				DiffSuppressFunc: helper.DurationDiffSuppress,
			},
		},
		CreateContext: update,
		UpdateContext: update,
		DeleteContext: delete,
		ReadContext:   read,
		Importer:      helper.ImportWithEmptyID(),
	}
}
