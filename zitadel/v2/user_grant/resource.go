package user_grant

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/helper"
)

func GetResource() *schema.Resource {
	return &schema.Resource{
		Description: "Resource representing the authorization given to a user directly, including the given roles.",
		Schema: map[string]*schema.Schema{
			helper.OrgIDVar: helper.OrgIDResourceField,
			ProjectIDVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID of the project",
				ForceNew:    true,
			},
			ProjectGrantIDVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID of the granted project",
				ForceNew:    true,
			},
			UserIDVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the user",
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
		},
		DeleteContext: delete,
		CreateContext: create,
		UpdateContext: update,
		ReadContext:   read,
		Importer: &schema.ResourceImporter{StateContext: helper.ImportWithIDAndOptionalOrgV5(
			helper.ResourceIDVar,
			helper.ImportAttribute{
				Key:             UserIDVar,
				ValueFromString: helper.ConvertID,
			},
			helper.ImportAttribute{
				Key:             ProjectIDVar,
				ValueFromString: helper.ConvertID,
				Optional:        true,
			},
			helper.ImportAttribute{
				Key:             ProjectGrantIDVar,
				ValueFromString: helper.ConvertID,
				Optional:        true,
			},
		)},
	}
}
