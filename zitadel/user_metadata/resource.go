package user_metadata

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/helper"
)

func GetResource() *schema.Resource {
	return &schema.Resource{
		Description: "Add a custom attribute to the user like the authenticating system. You can use this information in your actions. This Terraform resource manages a single key-value pair.",
		Schema: map[string]*schema.Schema{
			helper.OrgIDVar: helper.OrgIDResourceField,
			UserIDVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the user",
				ForceNew:    true,
			},
			KeyVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The key of a metadata entry",
				ForceNew:    true,
			},
			ValueVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The string representation of a metadata entry value. For binary data, use the base64encode function.",
			},
		},
		CreateContext: set,
		DeleteContext: delete,
		ReadContext:   read,
		UpdateContext: set,
		Importer:      helper.ImportWithEmptyID(helper.ImportOptionalOrgAttribute, helper.NewImportAttribute(UserIDVar, helper.ConvertID, false), helper.NewImportAttribute(KeyVar, helper.ConvertNonEmpty, false)),
	}
}
