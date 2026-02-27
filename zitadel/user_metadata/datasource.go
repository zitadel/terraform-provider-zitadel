package user_metadata

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
)

func GetDatasource() *schema.Resource {
	return &schema.Resource{
		Description: "Datasource representing a metadata entry of a user in ZITADEL.",
		Schema: map[string]*schema.Schema{
			helper.OrgIDVar: helper.OrgIDDatasourceField,
			UserIDVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the user",
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
		Description: "Datasource representing metadata entries of a user in ZITADEL.",
		Schema: map[string]*schema.Schema{
			helper.OrgIDVar: helper.OrgIDDatasourceField,
			UserIDVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the user",
			},
			KeyVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filter metadata by key",
			},
			metadataVar: {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of user metadata entries",
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
