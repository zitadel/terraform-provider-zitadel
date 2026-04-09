package action_target_public_key

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func GetDatasource() *schema.Resource {
	return &schema.Resource{
		Description: "Datasource representing a public key for an action target.",
		Schema: map[string]*schema.Schema{
			targetIDVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the action target.",
			},
			keyIDVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the public key.",
			},
			publicKeyVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The public key in PEM format (RSA or EC).",
			},
			activeVar: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the public key is active and used for payload encryption.",
			},
			fingerprintVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The fingerprint of the public key.",
			},
			expirationDateVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The expiration date of the public key in RFC3339 format.",
			},
			creationDateVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date the public key was added.",
			},
		},
		ReadContext: read,
	}
}
