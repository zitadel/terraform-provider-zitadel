package idp

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
		Description: "Datasource representing all identity providers on the instance, which can be looked up in detail with the type specific IDP datasources.",
		Schema: map[string]*schema.Schema{
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
		},
		ReadContext: list,
	}
}
