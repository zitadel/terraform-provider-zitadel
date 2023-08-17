package idp_gitlab_self_hosted

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

const IssuerVar = "issuer"

var (
	IssuerResourceField = &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "the providers issuer",
	}
	IssuerDataSourceField = &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "the providers issuer",
	}
)
