package idp_apple

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

const (
	TeamIDVar     = "team_id"
	KeyIDVar      = "key_id"
	PrivateKeyVar = "private_key"
)

var (
	TeamIDResourceField = &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "Apple Team ID from your Apple Developer Account",
	}
	TeamIDDataSourceField = &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Apple Team ID from your Apple Developer Account",
	}
	KeyIDResourceField = &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "Apple Key ID from your Apple Developer Account",
	}
	KeyIDDataSourceField = &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Apple Key ID from your Apple Developer Account",
	}
	PrivateKeyResourceField = &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "Apple Private Key from your Apple Developer Account",
		Sensitive:   true,
	}
	PrivateKeyDataSourceField = &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Apple Private Key from your Apple Developer Account",
		Sensitive:   true,
	}
)
