package instance_member

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
)

func GetResource() *schema.Resource {
	return &schema.Resource{
		Description: "Resource representing the membership of a user on an instance, defined with the given role.",
		Schema: map[string]*schema.Schema{
			UserIDVar: {
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
				Description: "List of roles granted. Instance member roles must start with 'IAM_' (e.g., IAM_OWNER, IAM_OWNER_VIEWER). See https://zitadel.com/docs/guides/manage/console/managers for available roles.",
			},
		},
		DeleteContext: delete,
		CreateContext: create,
		UpdateContext: update,
		ReadContext:   read,
		Importer:      helper.ImportWithEmptyID(helper.NewImportAttribute(UserIDVar, helper.ConvertID, false)),
	}
}
