package sms_provider_twilio

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
			tokenVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Token used to communicate with Twilio.",
				Sensitive:   true,
			},
			senderNumberVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Sender number which is used to send the SMS.",
			},
		},
		CreateContext: create,
		DeleteContext: delete,
		ReadContext:   read,
		UpdateContext: update,
		Importer:      &schema.ResourceImporter{StateContext: schema.ImportStatePassthroughContext},
	}
}
