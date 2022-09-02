package org

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

func GetDatasource() *schema.Resource {
	return &schema.Resource{
		Description: "Datasource representing an organization in ZITADEL, which is the highest level after the instance and contains several other resource including policies if the configuration differs to the default policies on the instance.",
		Schema: map[string]*schema.Schema{
			orgIDVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of this resource.",
			},
			nameVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the org",
			},
		},
		ReadContext: read,
		Importer:    &schema.ResourceImporter{StateContext: schema.ImportStatePassthroughContext},
	}
}
