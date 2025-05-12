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
				Required:    true,
				Description: "Description of the SMS provider.",
				Sensitive:   true,
			},
		},
		CreateContext: create,
		DeleteContext: delete,
		ReadContext:   read,
		UpdateContext: update,
		Importer:      helper.ImportWithID(IDVar),
	}
}
