package organization_metadata

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func GetDatasource() *schema.Resource {
	return &schema.Resource{
		Description: "Datasource representing a metadata entry of an organization in ZITADEL.",
		Schema: map[string]*schema.Schema{
			OrganizationIDVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the organization",
			},
			KeyVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Key of the metadata entry",
			},
			ValueVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Value of the metadata entry",
			},
		},
		ReadContext: get,
	}
}

func ListDatasources() *schema.Resource {
	return &schema.Resource{
		Description: "Datasource representing metadata entries of an organization in ZITADEL.",
		Schema: map[string]*schema.Schema{
			OrganizationIDVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the organization",
			},
			KeyVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filter metadata by key",
			},
			metadataVar: {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of organization metadata entries",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						KeyVar: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Key of the metadata entry",
						},
						ValueVar: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Value of the metadata entry",
						},
					},
				},
			},
		},
		ReadContext: list,
	}
}
