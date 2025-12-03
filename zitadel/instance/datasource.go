package instance

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func GetDatasource() *schema.Resource {
	return &schema.Resource{
		Description: "Datasource representing the ZITADEL instance with its domains.",
		Schema: map[string]*schema.Schema{
			InstanceIDVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID of the instance. If not provided, uses the instance from the authentication context.",
			},
			NameVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the instance.",
			},
			PrimaryDomainVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Primary domain of the instance. This is the first custom domain if any exist, otherwise the generated domain.",
			},
			GeneratedDomainVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The generated domain for this instance (e.g., instance1.zitadel.cloud).",
			},
			CustomDomainsVar: {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of custom domains configured for this instance.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			TrustedDomainsVar: {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of trusted domains configured for this instance.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
		ReadContext: list,
	}
}
