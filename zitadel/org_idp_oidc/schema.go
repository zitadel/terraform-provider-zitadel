package org_idp_oidc

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

const (
	IsIdTokenMappingVar = "is_id_token_mapping"
	IssuerVar           = "issuer"
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
		Description: "the issuer of the idp",
	}
	IssuerDatasourceField = &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "the issuer of the idp",
	}
)
