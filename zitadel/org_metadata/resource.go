package org_metadata

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
)

func GetResource() *schema.Resource {
	return &schema.Resource{
		Description: "Add a custom attribute to the organization like its location or an identifier in another system. You can use this information in your actions. This Terraform resource manages a single key-value pair.",
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
