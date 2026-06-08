package smtp_config

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
)

func GetResource() *schema.Resource {
	return &schema.Resource{
		Description:        "Instance SMTP configuration. **Deprecated:** the underlying SMTP config API is marked deprecated in ZITADEL. Use `zitadel_email_provider_smtp` instead.",
		DeprecationMessage: "The underlying SMTP config API is marked deprecated in ZITADEL. Use zitadel_email_provider_smtp instead.",
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
				WriteOnly:   true,
			},
			"password_hash": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "A non-reversible hash of the write-only password, used to detect when it changes. It does not contain the secret itself.",
			},
			replyToAddressVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Address to reply to.",
			},
			DescriptionVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of the SMTP configuration.",
			},
			SetActiveVar: {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Set the SMTP configuration active after creating/updating.",
			},
		},
		CreateContext: create,
		DeleteContext: delete,
		ReadContext:   read,
		UpdateContext: update,
		CustomizeDiff: func(ctx context.Context, d *schema.ResourceDiff, m interface{}) error {
			return helper.WriteOnlyHashDiff(d, PasswordVar, "password_hash")
		},
		Importer: helper.ImportWithIDAndOptionalSecret(IDVar, PasswordVar),
	}
}
