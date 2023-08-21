package org

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/helper"
)

func GetResource() *schema.Resource {
	return &schema.Resource{
		Description: "Resource representing an organization in ZITADEL, which is the highest level after the instance and contains several other resource including policies if the configuration differs to the default policies on the instance.",
		Schema: map[string]*schema.Schema{
			NameVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the org",
			},
			primaryDomainVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Primary domain of the org",
			},
			stateVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "State of the org",
			},
		},
		CreateContext: create,
		DeleteContext: delete,
		ReadContext:   get,
		UpdateContext: update,
		Importer:      helper.ImportWithID(OrgIDVar),
	}
}
