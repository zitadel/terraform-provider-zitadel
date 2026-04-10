package idp_saml

import (
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/idp"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
)

const (
	MetadataXMLVar                   = "metadata_xml"
	BindingVar                       = "binding"
	WithSignedRequestVar             = "with_signed_request"
	NameIdFormatVar                  = "name_id_format"
	TransientMappingAttributeNameVar = "transient_mapping_attribute_name"
	FederatedLogoutEnabledVar        = "federated_logout_enabled"
	SignatureAlgorithmVar            = "signature_algorithm"
)

var (
	BindingResourceField = &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "The binding" + helper.DescriptionEnumValuesList(idp.SAMLBinding_name),
		Default:     idp.SAMLBinding_name[0],
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
		Type:        schema.TypeBool,
		Computed:    true,
		Description: "Whether the SAML IDP requires signed requests",
	}
	NameIdFormatResourceField = &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "The nameid-format requested" + helper.DescriptionEnumValuesList(idp.SAMLNameIDFormat_name),
		Default:     idp.SAMLNameIDFormat_name[0],
		ValidateDiagFunc: func(value interface{}, path cty.Path) diag.Diagnostics {
			return helper.EnumValueValidation(NameIdFormatVar, value, idp.SAMLNameIDFormat_value)
		},
	}
	NameIdFormatDatasourceField = &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The nameid-format requested",
	}
	TransientMappingAttributeNameResourceField = &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Name of the attribute used to map the user in case the nameid-format is `urn:oasis:names:tc:SAML:2.0:nameid-format:transient`.",
	}
	TransientMappingAttributeNameDatasourceField = &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Name of the attribute used to map the user in case the nameid-format is `urn:oasis:names:tc:SAML:2.0:nameid-format:transient`.",
	}
	FederatedLogoutEnabledResourceField = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "If enabled, ZITADEL will send a logout request to the identity provider when the user terminates the session in ZITADEL. Be sure to provide a SLO endpoint as part of the metadata.",
	}
	FederatedLogoutEnabledDatasourceField = &schema.Schema{
		Type:        schema.TypeBool,
		Computed:    true,
		Description: "If enabled, ZITADEL will send a logout request to the identity provider when the user terminates the session in ZITADEL.",
	}
	SignatureAlgorithmResourceField = &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Signature Algorithm used to sign SAML requests and responses. Can be used only if `with_signed_request` is true." + helper.DescriptionEnumValuesList(idp.SAMLSignatureAlgorithm_name),
		Default:     idp.SAMLSignatureAlgorithm_name[0],
		ValidateDiagFunc: func(value interface{}, path cty.Path) diag.Diagnostics {
			return helper.EnumValueValidation(SignatureAlgorithmVar, value, idp.SAMLSignatureAlgorithm_value)
		},
	}
	SignatureAlgorithmDatasourceField = &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Signature Algorithm used to sign SAML requests and responses.",
	}
)
