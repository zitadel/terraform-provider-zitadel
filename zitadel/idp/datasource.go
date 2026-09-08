package idp

import (
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/idp"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/object"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/idp_utils"
)

func ListDatasources() *schema.Resource {
	return &schema.Resource{
		Description: "Datasource representing all identity providers on the instance, optionally filtered by name and type.",
		Schema: map[string]*schema.Schema{
			idp_utils.NameVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name to filter identity providers by",
			},
			NameMethodVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Method for querying identity providers by name" + helper.DescriptionEnumValuesList(object.TextQueryMethod_name),
				ValidateDiagFunc: func(value interface{}, path cty.Path) diag.Diagnostics {
					return helper.EnumValueValidation(NameMethodVar, value, object.TextQueryMethod_value)
				},
				Default: object.TextQueryMethod_TEXT_QUERY_METHOD_EQUALS_IGNORE_CASE.String(),
			},
			TypeVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Type to filter identity providers by, PROVIDER_TYPE_UNSPECIFIED applies no filter" + helper.DescriptionEnumValuesList(idp.ProviderType_name),
				ValidateDiagFunc: func(value interface{}, path cty.Path) diag.Diagnostics {
					return helper.EnumValueValidation(TypeVar, value, idp.ProviderType_value)
				},
			},
			idpsVar: {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of identity providers",
				Elem:        IdpElem(),
			},
		},
		ReadContext: list,
	}
}

// IdpElem is the nested schema of a listed identity provider, shared by the instance and organization datasources.
func IdpElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			idp_utils.IdpIDVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the identity provider",
			},
			idp_utils.NameVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the identity provider",
			},
			TypeVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Type of the identity provider",
			},
			stateVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "State of the identity provider",
			},
			ownerTypeVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Owner type of the identity provider, either the instance (system) or an organization",
			},
		},
	}
}
