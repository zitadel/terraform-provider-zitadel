package smtp_config

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/helper"
)

func GetResource() *schema.Resource {
	return &schema.Resource{
		Description: "Resource representing the SMTP configuration of an instance.",
		Schema: map[string]*schema.Schema{
			SenderAddressVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Address used to send emails.",
			},
			SenderNameVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Sender name used to send emails.",
			},
			tlsVar: {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "TLS used to communicate with your SMTP server.",
			},
			hostVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Host and port address to your SMTP server.",
			},
			userVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "User used to communicate with your SMTP server.",
			},
			PasswordVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Password used to communicate with your SMTP server.",
				Sensitive:   true,
			},
			replyToAddressVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Address to reply to.",
			},
			SetActiveVar: {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Set the SMTP configuration active after creating/updating",
			},
		},
		CreateContext: create,
		DeleteContext: delete,
		ReadContext:   read,
		UpdateContext: update,
		Importer:      helper.ImportWithEmptyID(helper.NewImportAttribute(PasswordVar, helper.ConvertNonEmpty, true)),
	}
}
