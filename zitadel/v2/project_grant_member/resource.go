package project_grant_member

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/helper"
)

func GetResource() *schema.Resource {
	return &schema.Resource{
		Description: "Resource representing the membership of a user on an granted project, defined with the given role.",
		Schema: map[string]*schema.Schema{
			helper.OrgIDVar: helper.OrgIDResourceField,
			projectIDVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the project",
				ForceNew:    true,
			},
			GrantIDVar: {
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
			RolesVar: {
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
		Importer: helper.ImportWithEmptyID(
			helper.ImportOptionalOrgAttribute,
			helper.NewImportAttribute(projectIDVar, helper.ConvertID, false),
			helper.NewImportAttribute(GrantIDVar, helper.ConvertID, false),
			helper.NewImportAttribute(userIDVar, helper.ConvertID, false),
		),
	}
}
