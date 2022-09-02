package project_grant_member

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func GetResource() *schema.Resource {
	return &schema.Resource{
		Description: "Resource representing the membership of a user on an granted project, defined with the given role.",
		Schema: map[string]*schema.Schema{
			orgIDVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the organization which owns the resource",
				ForceNew:    true,
			},
			projectIDVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the project",
				ForceNew:    true,
			},
			grantIDVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the grant",
				ForceNew:    true,
			},
			userIDVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the user",
				ForceNew:    true,
			},
			rolesVar: {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Required:    true,
				Description: "List of roles granted",
			},
		},
		DeleteContext: delete,
		CreateContext: create,
		UpdateContext: update,
		ReadContext:   read,
		Importer:      &schema.ResourceImporter{StateContext: schema.ImportStatePassthroughContext},
	}
}
