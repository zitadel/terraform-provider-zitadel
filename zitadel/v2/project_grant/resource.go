package project_grant

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func GetResource() *schema.Resource {
	return &schema.Resource{
		Description: "Resource representing the grant of a project to a different organization, also containing the available roles which can be given to the members of the projectgrant.",
		Schema: map[string]*schema.Schema{
			projectIDVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the project",
				ForceNew:    true,
			},
			grantedOrgIDVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the organization granted the project",
				ForceNew:    true,
			},
			roleKeysVar: {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				Description: "List of roles granted",
			},
			orgIDVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the organization which owns the resource",
			},
		},
		DeleteContext: delete,
		CreateContext: create,
		UpdateContext: update,
		ReadContext:   read,
		Importer:      &schema.ResourceImporter{StateContext: schema.ImportStatePassthroughContext},
	}
}
