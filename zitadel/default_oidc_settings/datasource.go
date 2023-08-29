package default_oidc_settings

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func GetDatasource() *schema.Resource {
	return &schema.Resource{
		Description: "Datasource representing the default oidc settings.",
		Schema: map[string]*schema.Schema{
			accessTokenLifetimeVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "lifetime duration of access tokens",
			},
			idTokenLifetimeVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "lifetime duration of id tokens",
			},
			RefreshTokenExpirationVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "expiration duration of refresh tokens",
			},
			refreshTokenIdleExpirationVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "expiration duration of idle refresh tokens",
			},
		},
		ReadContext: read,
	}
}
