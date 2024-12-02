package org_idp_jwt

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
)

func GetDatasource() *schema.Resource {
	return &schema.Resource{
		Description: "Datasource representing a generic JWT IdP on the organization.",
		Schema: map[string]*schema.Schema{
			idpIDVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of this resource.",
			},
			helper.OrgIDVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the organization",
			},
			nameVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the IDP",
			},
			stylingTypeVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Some identity providers specify the styling of the button to their login",
			},
			JwtEndpointVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the endpoint where the jwt can be extracted",
			},
			keysEndpointVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the endpoint to the key (JWK) which are used to sign the JWT with",
			},
			issuerVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the issuer of the jwt (for validation)",
			},
			headerNameVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the name of the header where the JWT is sent in, default is authorization",
			},
			autoRegisterVar: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "auto register for users from this idp",
			},
		},
		ReadContext: read,
	}
}
