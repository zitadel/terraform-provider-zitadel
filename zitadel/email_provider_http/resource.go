package email_provider_http

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
)

func GetResource() *schema.Resource {
	return &schema.Resource{
		Description: "Resource representing the HTTP email provider configuration of an instance.",
		Schema: map[string]*schema.Schema{
			EndpointVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Http endpoint which is used to send the email.",
			},
			DescriptionVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of the email provider.",
			},
			SigningKeyVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "Key used to sign and check payload sent to the HTTP provider.",
			},
			setActiveVar: {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Set the email provider as active after creating/updating.",
			},
		},
		CreateContext: create,
		DeleteContext: delete,
		ReadContext:   read,
		UpdateContext: update,
		Importer:      helper.ImportWithIDAndOptionalSecret(IDVar, SigningKeyVar),
	}
}
