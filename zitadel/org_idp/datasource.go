package org_idp

import (
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/idp"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/object"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
)

func ListDatasources() *schema.Resource {
	return &schema.Resource{
		Description: "Datasource representing all identity providers available to an organization, which can be looked up in detail with the type specific IdP datasources.",
		Schema: map[string]*schema.Schema{
			helper.OrgIDVar: helper.OrgIDDatasourceField,
			idpIDsVar: {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of all identity provider IDs.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			NameVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the identity provider.",
			},
			nameMethodVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Method for querying identity providers by name" + helper.DescriptionEnumValuesList(object.TextQueryMethod_name),
				ValidateDiagFunc: func(value interface{}, path cty.Path) diag.Diagnostics {
					return helper.EnumValueValidation(nameMethodVar, value, object.TextQueryMethod_value)
				},
				Default: object.TextQueryMethod_TEXT_QUERY_METHOD_EQUALS_IGNORE_CASE.String(),
			},
			typeVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Type of the identity provider" + helper.DescriptionEnumValuesList(idp.ProviderType_name),
				ValidateDiagFunc: func(value interface{}, path cty.Path) diag.Diagnostics {
					return helper.EnumValueValidation(typeVar, value, idp.ProviderType_value)
				},
			},
			ownerTypeVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Owner type of the identity provider, either the instance (system) or the organization" + helper.DescriptionEnumValuesList(idp.IDPOwnerType_name),
				ValidateDiagFunc: func(value interface{}, path cty.Path) diag.Diagnostics {
					return helper.EnumValueValidation(ownerTypeVar, value, idp.IDPOwnerType_value)
				},
			},
		},
		ReadContext: list,
	}
}
