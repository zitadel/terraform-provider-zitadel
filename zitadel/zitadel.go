package zitadel

import (
	"context"
	"errors"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
)

func GetZitadelDatasource() *schema.Resource {
	return &schema.Resource{
		Description: "Datasource representing the session token of the provider's configuration.",
		Schema: map[string]*schema.Schema{
			"token": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "The session token.",
			},
		},
		ReadContext: func(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
			clientInfo, ok := m.(*helper.ClientInfo)
			if !ok {
				return diag.Errorf("failed to get client info")
			}

			tokenSource, err := clientInfo.JWTProfileTokenSource(ctx)
			if errors.Is(err, helper.ErrNoJWTProfile) {
				return diag.Errorf("Session token generation is only supported when using 'jwt_profile_file', 'jwt_profile_json' or 'token' (service account key file path) in the provider configuration.")
			}
			if err != nil {
				return diag.FromErr(err)
			}

			token, err := tokenSource.Token()
			if err != nil {
				return diag.FromErr(err)
			}

			// Data sources must have an ID set. We use the current time to ensure uniqueness
			// and indicate when it was read.
			d.SetId(time.Now().UTC().Format(time.RFC3339Nano))

			if err := d.Set("token", token.AccessToken); err != nil {
				return diag.FromErr(err)
			}

			return nil
		},
	}
}
