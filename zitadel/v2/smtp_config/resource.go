package smtp_config

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func GetResource() *schema.Resource {
	return &schema.Resource{
		Description: "Resource representing the SMTP configuration of an instance.",
		Schema: map[string]*schema.Schema{
			senderAddressVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Address used to send emails.",
			},
			senderNameVar: {
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
			passwordVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Password used to communicate with your SMTP server.",
				Sensitive:   true,
			},
		},
		CreateContext: create,
		DeleteContext: delete,
		ReadContext:   read,
		UpdateContext: update,
		Importer:      &schema.ResourceImporter{StateContext: schema.ImportStatePassthroughContext},
	}
}
