package org

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
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
			IsDefaultVar: {
				Type:                  schema.TypeBool,
				Optional:              true,
				DiffSuppressOnRefresh: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					// Need to avoid forever pending changes when trying to set this to false
					// since setting to false will not actually unmark the org as default.
					return new != "true"
				},
				Description: "True sets the org as default org for the instance. Only one org can be default org. Nothing happens if you set it to false until you set another org as default org.",
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
			OrgIDInputVar: {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Optionally set a custom unique ID for the organization. If omitted, ZITADEL will generate one.",
			},
		},
		CreateContext: create,
		DeleteContext: delete,
		ReadContext:   get,
		UpdateContext: update,
		Importer:      helper.ImportWithID(OrgIDVar),
	}
}
