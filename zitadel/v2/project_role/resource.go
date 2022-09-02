package project_role

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func GetResource() *schema.Resource {
	return &schema.Resource{
		Description: "Resource representing the project roles, which can be given as authorizations to users.",
		Schema: map[string]*schema.Schema{
			projectIDVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the project",
				ForceNew:    true,
			},
			orgIDVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the organization",
				ForceNew:    true,
			},
			keyVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Key used for project role",
			},
			displayNameVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name used for project role",
			},
			groupVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Group used for project role",
			},
		},
		DeleteContext: delete,
		CreateContext: create,
		UpdateContext: update,
		ReadContext:   read,
		Importer:      &schema.ResourceImporter{StateContext: schema.ImportStatePassthroughContext},
	}
}
