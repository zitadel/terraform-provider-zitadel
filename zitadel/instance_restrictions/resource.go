package instance_restrictions

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
)

func GetResource() *schema.Resource {
	return &schema.Resource{
		Description: "Resource representing the instance restrictions.",
		Schema: map[string]*schema.Schema{
			disallowPublicOrgRegistrationVar: {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Disallow public organization registration for the instance",
			},
			allowedLanguagesVar: {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				Computed:    true,
				Description: "Allowed languages (BCP 47 language tags) for the instance",
			},
		},
		CreateContext: create,
		ReadContext:   read,
		UpdateContext: update,
		DeleteContext: delete,
		Importer:      helper.ImportWithEmptyID(),
	}
}
