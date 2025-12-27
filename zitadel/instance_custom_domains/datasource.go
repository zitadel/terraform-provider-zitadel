package instance_custom_domains

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func GetDatasource() *schema.Resource {
	return &schema.Resource{
		Description: "Datasource representing the list of custom domains configured for a ZITADEL instance.",
		Schema: map[string]*schema.Schema{
			InstanceIDVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ID of the instance. If not provided, the instance from the current context (e.g., identified by the host header) will be used.",
			},
			DomainsVar: {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of custom domain names",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
		ReadContext: list,
	}
}
