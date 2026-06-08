package idp_apple

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

const (
	TeamIDVar         = "team_id"
	KeyIDVar          = "key_id"
	PrivateKeyVar     = "private_key"
	PrivateKeyHashVar = "private_key_hash"
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
		WriteOnly:   true,
	}
	PrivateKeyHashResourceField = &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Sensitive:   true,
		Description: "A non-reversible hash of the write-only private_key, used to detect when it changes. It does not contain the key itself.",
	}
)
