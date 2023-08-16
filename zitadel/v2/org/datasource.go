package org

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/helper"
)

func GetDatasource() *schema.Resource {
	return &schema.Resource{
		Description: "Datasource representing an organization in ZITADEL, which is the highest level after the instance and contains several other resource including policies if the configuration differs to the default policies on the instance.",
		Schema: map[string]*schema.Schema{
			helper.ResourceIDVar: helper.ResourceIDDatasourceField,
			nameVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the org",
			},
		},
		ReadContext: read,
	}
}
