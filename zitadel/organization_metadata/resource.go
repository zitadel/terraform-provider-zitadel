package organization_metadata

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
)

func GetResource() *schema.Resource {
	return &schema.Resource{
		Description: "Resource representing metadata of an organization in ZITADEL. This resource manages a single key-value pair.",
		Schema: map[string]*schema.Schema{
			OrganizationIDVar: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the organization",
			},
			KeyVar: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Key of the metadata entry",
			},
			ValueVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Value of the metadata entry. For binary data, use base64encode function.",
			},
		},
		CreateContext: set,
		DeleteContext: delete,
		ReadContext:   read,
		UpdateContext: set,
		Importer:      helper.ImportWithOptionalOrg(helper.NewImportAttribute(KeyVar, helper.ConvertNonEmpty, false)),
	}
}
