package org_idp

import (
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	idppb "github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/idp"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/object"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/idp"
)

func ListDatasources() *schema.Resource {
	return &schema.Resource{
		Description: "Datasource representing all identity providers available to an organization, optionally filtered by name, type and owner.",
		Schema: map[string]*schema.Schema{
			helper.OrgIDVar: helper.OrgIDDatasourceField,
			idp.NameVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name to filter identity providers by",
			},
			"name_method": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Method for querying identity providers by name" + helper.DescriptionEnumValuesList(object.TextQueryMethod_name),
				ValidateDiagFunc: func(value interface{}, path cty.Path) diag.Diagnostics {
					return helper.EnumValueValidation("name_method", value, object.TextQueryMethod_value)
				},
				Default: object.TextQueryMethod_TEXT_QUERY_METHOD_EQUALS_IGNORE_CASE.String(),
			},
			idp.TypeVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Type to filter identity providers by" + helper.DescriptionEnumValuesList(idppb.ProviderType_name),
				ValidateDiagFunc: func(value interface{}, path cty.Path) diag.Diagnostics {
					return helper.EnumValueValidation(idp.TypeVar, value, idppb.ProviderType_value)
				},
			},
			ownerTypeVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Owner type to filter identity providers by, either the instance (system) or the organization" + helper.DescriptionEnumValuesList(idppb.IDPOwnerType_name),
				ValidateDiagFunc: func(value interface{}, path cty.Path) diag.Diagnostics {
					return helper.EnumValueValidation(ownerTypeVar, value, idppb.IDPOwnerType_value)
				},
			},
			idpsVar: {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of identity providers",
				Elem:        idp.IdpElem(),
			},
		},
		ReadContext: list,
	}
}
