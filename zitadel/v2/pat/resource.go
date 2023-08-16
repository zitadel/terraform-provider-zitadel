package pat

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/helper"
)

func GetResource() *schema.Resource {
	return &schema.Resource{
		Description: "Resource representing a personal access token of a user",
		Schema: map[string]*schema.Schema{
			helper.OrgIDVar: helper.OrgIDResourceField,
			UserIDVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the user",
				ForceNew:    true,
			},
			TokenVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Value of the token",
				Sensitive:   true,
			},
			expirationDateVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Expiration date of the token in the RFC3339 format",
				ForceNew:    true,
			},
		},
		DeleteContext: delete,
		CreateContext: create,
		ReadContext:   read,
		Importer: &schema.ResourceImporter{StateContext: helper.ImportWithIDAndOptionalOrg(
			helper.ResourceIDVar,
			helper.NewImportAttribute(UserIDVar, helper.ConvertID, false),
			helper.NewImportAttribute(TokenVar, helper.ConvertNonEmpty, true),
		)},
	}
}
