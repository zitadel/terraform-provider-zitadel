package org_metadata

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
)

func GetResource() *schema.Resource {
	return &schema.Resource{
		Description:        "Custom key-value metadata on an organization. **Deprecated:** use `zitadel_organization_metadata` which uses the metadata/v2 API (requires ZITADEL 4.x).",
		DeprecationMessage: "Use zitadel_organization_metadata instead (metadata/v2 API, requires ZITADEL 4.x).",
		Schema: map[string]*schema.Schema{
			helper.OrgIDVar: helper.OrgIDResourceField,
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
		Importer:      helper.ImportWithOptionalOrg(helper.NewImportAttribute(KeyVar, helper.ConvertNonEmpty, false)),
	}
}
