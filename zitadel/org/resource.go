package org

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
)

func GetResource() *schema.Resource {
	return &schema.Resource{
		Description: "Organization resource that is backward-compatible with ZITADEL 3.x and 4.x: it tries the org/v2 API first and falls back to the management API automatically. For deployments that only target ZITADEL 4.x, prefer `zitadel_organization`.",
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
			adminsVar: {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "ID of existing user to grant admin access",
						},
						"roles": {
							Type:        schema.TypeList,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "Roles to assign (defaults to ORG_OWNER if empty)",
						},
					},
				},
				Description: "Admin users for the organization",
			},
		},
		CreateContext: create,
		DeleteContext: delete,
		ReadContext:   get,
		UpdateContext: update,
		Importer:      helper.ImportWithID(OrgIDVar),
	}
}
