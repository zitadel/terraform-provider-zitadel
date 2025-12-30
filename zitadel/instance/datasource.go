package instance

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func GetDatasource() *schema.Resource {
	return &schema.Resource{
		Description: "Datasource representing the ZITADEL instance.",
		Schema: map[string]*schema.Schema{
			InstanceIDVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the instance.",
			},
			NameVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the instance.",
			},
			VersionVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Version of the ZITADEL system the instance is running on.",
			},
			StateVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Current state of the instance.",
			},
			CustomDomainsVar: {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Custom domains configured for this instance.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"domain": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The domain name.",
						},
						"primary": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether this is the primary domain.",
						},
						"generated": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether this domain was auto-generated.",
						},
					},
				},
			},
			TrustedDomainsVar: {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Trusted domains configured for this instance.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
		ReadContext: list,
	}
}
