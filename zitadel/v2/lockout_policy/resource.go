package lockout_policy

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func GetResource() *schema.Resource {
	return &schema.Resource{
		Description: "Resource representing the custom lockout policy of an organization.",
		Schema: map[string]*schema.Schema{
			orgIDVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Id for the organization",
				ForceNew:    true,
			},
			maxPasswordAttemptsVar: {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Maximum password check attempts before the account gets locked. Attempts are reset as soon as the password is entered correct or the password is reset.",
			},
		},
		DeleteContext: delete,
		CreateContext: create,
		UpdateContext: update,
		ReadContext:   read,
		Importer:      &schema.ResourceImporter{StateContext: schema.ImportStatePassthroughContext},
	}
}
