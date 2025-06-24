package webkey

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
)

func GetDatasource() *schema.Resource {
	return &schema.Resource{
		Description: "Datasource representing a web key.",
		Schema: map[string]*schema.Schema{
			WebKeyIDVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of this resource.",
			},
			helper.OrgIDVar: helper.OrgIDResourceField,
			StateVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "State of the key.",
			},
			KeyTypeVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Type of the key.",
			},
			PubKeyVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "The public key, PEM encoded.",
			},
		},
		ReadContext: readWebKey,
	}
}
