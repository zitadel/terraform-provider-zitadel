package pat

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func GetResource() *schema.Resource {
	return &schema.Resource{
		Description: "Resource representing a personal access token of a user",
		Schema: map[string]*schema.Schema{
			orgIDVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the organization",
				ForceNew:    true,
			},
			userIDVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the user",
				ForceNew:    true,
			},
			tokenVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Value of the token",
				Sensitive:   true,
			},
			expirationDateVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Expiration date of the token",
				ForceNew:    true,
			},
		},
		DeleteContext: delete,
		CreateContext: create,
		ReadContext:   read,
		Importer:      &schema.ResourceImporter{StateContext: schema.ImportStatePassthroughContext},
	}
}
