package project_role

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
)

func GetResource() *schema.Resource {
	return &schema.Resource{
		Description: "Resource representing the project roles, which can be given as authorizations to users.",
		Schema: map[string]*schema.Schema{
			helper.OrgIDVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the organization",
			},
			ProjectIDVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the project",
				ForceNew:    true,
			},
			KeyVar: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
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
		Importer: helper.ImportWithEmptyID(
			helper.NewImportAttribute(ProjectIDVar, helper.ConvertID, false),
			helper.NewImportAttribute(KeyVar, helper.ConvertNonEmpty, false),
			helper.ImportOptionalOrgAttribute,
		),
	}
}
