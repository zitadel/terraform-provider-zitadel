package idp_saml

import (
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/idp"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
)

const (
	MetadataXMLVar       = "metadata_xml"
	BindingVar           = "binding"
	WithSignedRequestVar = "with_signed_request"
)

var (
	BindingResourceField = &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "The binding" + helper.DescriptionEnumValuesList(idp.SAMLBinding_name),
		Default:     idp.SAMLBinding_SAML_BINDING_UNSPECIFIED,
		ValidateDiagFunc: func(value interface{}, path cty.Path) diag.Diagnostics {
			return helper.EnumValueValidation(BindingVar, value, idp.SAMLBinding_value)
		},
	}
	BindingDatasourceField = &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The binding",
	}
	MetadataXMLResourceField = &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "The metadata XML as plain string",
	}
	MetadataXMLDatasourceField = &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The metadata XML as plain string",
	}
	WithSignedRequestResourceField = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Whether the SAML IDP requires signed requests",
	}
	WithSignedRequestDatasourceField = &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Whether the SAML IDP requires signed requests",
	}
)
