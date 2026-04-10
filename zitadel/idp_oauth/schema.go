package idp_oauth

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	AuthorizationEndpointVar = "authorization_endpoint"
	TokenEndpointVar         = "token_endpoint"
	UserEndpointVar          = "user_endpoint"
	IdAttributeVar           = "id_attribute"
	UsePKCEVar               = "use_pkce"
)

var (
	AuthorizationEndpointResourceField = &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "The authorization endpoint",
	}
	AuthorizationEndpointDatasourceField = &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The authorization endpoint",
	}
	TokenEndpointResourceField = &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "The token endpoint",
	}
	TokenEndpointDatasourceField = &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The token endpoint",
	}
	UserEndpointResourceField = &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "The user endpoint",
	}
	UserEndpointDatasourceField = &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The user endpoint",
	}
	IdAttributeResourceField = &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "The id attribute",
	}
	IdAttributeDatasourceField = &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The id attribute",
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
