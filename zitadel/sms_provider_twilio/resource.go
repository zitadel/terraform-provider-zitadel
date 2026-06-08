package sms_provider_twilio

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
)

func GetResource() *schema.Resource {
	return &schema.Resource{
		Description: "Resource representing the SMS provider Twilio configuration of an instance.",
		Schema: map[string]*schema.Schema{
			sidVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "SID used to communicate with Twilio.",
			},
			TokenVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Token used to communicate with Twilio.",
				Sensitive:   true,
				WriteOnly:   true,
			},
			"token_hash": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "A non-reversible hash of the write-only token, used to detect when it changes. It does not contain the secret itself.",
			},
			SenderNumberVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Sender number which is used to send the SMS.",
			},
			setActiveVar: {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Set the SMS provider as active after creating/updating.",
			},
			VerifyServiceSidVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Twilio Verify Service SID used for phone verification.",
			},
			DescriptionVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of the SMS provider.",
			},
		},
		CreateContext: create,
		DeleteContext: delete,
		ReadContext:   read,
		UpdateContext: update,
		CustomizeDiff: func(ctx context.Context, d *schema.ResourceDiff, m interface{}) error {
			return helper.WriteOnlyHashDiff(d, TokenVar, "token_hash")
		},
		Importer: helper.ImportWithIDAndOptionalSecret(providerIDVar, TokenVar),
	}
}
