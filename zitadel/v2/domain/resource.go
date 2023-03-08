package domain

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func GetResource() *schema.Resource {
	return &schema.Resource{
		Description: "Resource representing a domain of the organization.",
		Schema: map[string]*schema.Schema{
			nameVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the domain",
				ForceNew:    true,
			},
			orgIDVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the organization",
				ForceNew:    true,
			},
			isVerifiedVar: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Is domain verified",
			},
			isPrimaryVar: {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Is domain primary",
			},
			validationTypeVar: {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Validation type",
			},
		},
		ReadContext:   read,
		CreateContext: create,
		UpdateContext: update,
		DeleteContext: delete,
		Importer:      &schema.ResourceImporter{StateContext: schema.ImportStatePassthroughContext},
	}
}
