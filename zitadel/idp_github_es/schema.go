package idp_github_es

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

const (
	AuthorizationEndpointVar = "authorization_endpoint"
	TokenEndpointVar         = "token_endpoint"
	UserEndpointVar          = "user_endpoint"
)

var (
	AuthorizationEndpointResourceField = &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "the providers authorization endpoint",
	}
	AuthorizationEndpointDataSourceField = &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "the providers authorization endpoint",
	}
	TokenEndpointResourceField = &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "the providers token endpoint",
	}
	TokenEndpointDataSourceField = &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "the providers token endpoint",
	}
	UserEndpointResourceField = &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "the providers user endpoint",
	}
	UserEndpointDataSourceField = &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "the providers user endpoint",
	}
)
