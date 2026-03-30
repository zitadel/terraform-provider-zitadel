package action_target_public_key

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
)

func GetResource() *schema.Resource {
	return &schema.Resource{
		Description: "Resource representing a public key for an action target, used for payload encryption.",
		Schema: map[string]*schema.Schema{
			targetIDVar: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the action target to add the public key to.",
			},
			publicKeyVar: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The public key in PEM format (RSA or EC).",
			},
			expirationDateVar: {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The expiration date of the public key in RFC3339 format.",
			},
			keyIDVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the public key, used as 'kid' in the JWE header.",
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
			creationDateVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date the public key was added.",
			},
		},
		CreateContext: create,
		ReadContext:   read,
		DeleteContext: delete,
		Importer: helper.ImportWithID(
			keyIDVar,
			helper.NewImportAttribute(targetIDVar, helper.ConvertID, false),
		),
	}
}
