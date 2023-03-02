package default_notification_policy

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func GetResource() *schema.Resource {
	return &schema.Resource{
		Description: "Resource representing the default notification policy.",
		Schema: map[string]*schema.Schema{
			passwordChangeVar: {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Send notification if a user changes his password",
			},
		},
		ReadContext:   read,
		CreateContext: update,
		DeleteContext: delete,
		UpdateContext: update,
		Importer:      &schema.ResourceImporter{StateContext: schema.ImportStatePassthroughContext},
	}
}
