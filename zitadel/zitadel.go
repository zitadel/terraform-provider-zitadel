package zitadel

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zitadel/oidc/v3/pkg/client/profile"
	"github.com/zitadel/oidc/v3/pkg/oidc"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel"
	"golang.org/x/oauth2"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
)

func GetSessionTokenDatasource() *schema.Resource {
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

			scopes := []string{oidc.ScopeOpenID, zitadel.ScopeZitadelAPI()}
			var tokenSource oauth2.TokenSource
			var err error

			if clientInfo.KeyPath != "" {
				tokenSource, err = profile.NewJWTProfileTokenSourceFromKeyFile(ctx, clientInfo.Issuer, clientInfo.KeyPath, scopes)
			} else if len(clientInfo.Data) > 0 {
				tokenSource, err = profile.NewJWTProfileTokenSourceFromKeyFileData(ctx, clientInfo.Issuer, clientInfo.Data, scopes)
			} else {
				return diag.Errorf("Session token generation is only supported when using 'jwt_profile_file', 'jwt_profile_json' or 'token' (service account key) in the provider configuration.")
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
