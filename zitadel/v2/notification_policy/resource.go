package notification_policy

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func GetResource() *schema.Resource {
	return &schema.Resource{
		Description: "Resource representing the custom notification policy of an organization.",
		Schema: map[string]*schema.Schema{
			orgIDVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Id for the organization",
				ForceNew:    true,
			},
			passwordChangeVar: {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Send notification if a user changes his password",
			},
		},
		ReadContext:   read,
		CreateContext: create,
		DeleteContext: delete,
		UpdateContext: update,
		Importer:      &schema.ResourceImporter{StateContext: schema.ImportStatePassthroughContext},
	}
}
