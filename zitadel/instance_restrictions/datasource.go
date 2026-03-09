package instance_restrictions

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func GetDatasource() *schema.Resource {
	return &schema.Resource{
		Description: "Datasource representing the instance restrictions.",
		Schema: map[string]*schema.Schema{
			disallowPublicOrgRegistrationVar: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Disallow public organization registration for the instance",
			},
			allowedLanguagesVar: {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Allowed languages (BCP 47 language tags) for the instance",
			},
		},
		ReadContext: read,
	}
}
