package sms_provider_twilio

import (
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
		},
		CreateContext: create,
		DeleteContext: delete,
		ReadContext:   read,
		UpdateContext: update,
		Importer:      helper.ImportWithIDAndOptionalSecret(providerIDVar, TokenVar),
	}
}
