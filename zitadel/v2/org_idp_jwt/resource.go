package org_idp_jwt

import (
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/idp"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/helper"
)

func GetResource() *schema.Resource {
	return &schema.Resource{
		Description: "Resource representing a domain of the organization.",
		Schema: map[string]*schema.Schema{
			orgIDVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the organization",
				ForceNew:    true,
			},
			nameVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the IDP",
			},
			stylingTypeVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Some identity providers specify the styling of the button to their login" + helper.DescriptionEnumValuesList(idp.IDPStylingType_name),
				ValidateDiagFunc: func(value interface{}, path cty.Path) diag.Diagnostics {
					return helper.EnumValueValidation(stylingTypeVar, value, idp.IDPStylingType_value)
				},
			},
			jwtEndpointVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "the endpoint where the jwt can be extracted",
			},
			keysEndpointVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "the endpoint to the key (JWK) which are used to sign the JWT with",
			},
			issuerVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "the issuer of the jwt (for validation)",
			},
			headerNameVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "the name of the header where the JWT is sent in, default is authorization",
			},
			autoRegisterVar: {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "auto register for users from this idp",
			},
		},
		ReadContext:   read,
		CreateContext: create,
		UpdateContext: update,
		DeleteContext: delete,
		Importer:      &schema.ResourceImporter{StateContext: schema.ImportStatePassthroughContext},
	}
}
