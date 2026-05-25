package organization_metadata

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
)

func GetResource() *schema.Resource {
	return &schema.Resource{
		Description: "Custom key-value metadata on an organization, using the metadata/v2 API. **Requires ZITADEL 4.x.** For 3.x compatibility use `zitadel_org_metadata`.",
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
