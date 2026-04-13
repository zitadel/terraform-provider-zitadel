package sms_provider_http

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
)

func GetResource() *schema.Resource {
	return &schema.Resource{
		Description: "Resource representing the HTTP SMS provider configuration of an instance.",
		Schema: map[string]*schema.Schema{
			EndPointVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Http endpoint which is used to send the SMS.",
			},
			DescriptionVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of the SMS provider.",
			},
			setActiveVar: {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Set the SMS provider as active after creating/updating.",
			},
			SigningKeyVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "Key used to sign and check payload sent to the HTTP provider",
			},
			ExpirationSigningKeyVar: {
				Type:             schema.TypeString,
				Optional:         true,
				Description:      "Expiration duration for the signing key. When set during update, the old signing key will remain valid for the specified duration to allow for a graceful key rotation.",
				DiffSuppressFunc: helper.DurationDiffSuppress,
			},
		},
		CreateContext: create,
		DeleteContext: delete,
		ReadContext:   read,
		UpdateContext: update,
		Importer:      helper.ImportWithIDAndOptionalSecret(IDVar, SigningKeyVar),
	}
}
