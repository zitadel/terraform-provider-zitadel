package idp_oidc

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

const (
	IsIdTokenMappingVar = "is_id_token_mapping"
	IssuerVar           = "issuer"
	UsePKCEVar          = "use_pkce"
)

var (
	IsIdTokenMappingResourceField = &schema.Schema{
		Type:        schema.TypeBool,
		Required:    true,
		Description: "if true, provider information get mapped from the id token, not from the userinfo endpoint",
	}
	IsIdTokenMappingDatasourceField = &schema.Schema{
		Type:        schema.TypeBool,
		Computed:    true,
		Description: "if true, provider information get mapped from the id token, not from the userinfo endpoint.",
	}
	IssuerResourceField = &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "the OIDC issuer of the identity provider",
	}
	IssuerDatasourceField = &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "the OIDC issuer of the identity provider",
	}
	UsePKCEResourceField = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Computed:    true,
		Description: "Defines if the Proof Key for Code Exchange (PKCE) is used for the authorization code flow.",
	}
	UsePKCEDatasourceField = &schema.Schema{
		Type:        schema.TypeBool,
		Computed:    true,
		Description: "Defines if the Proof Key for Code Exchange (PKCE) is used for the authorization code flow.",
	}
)
