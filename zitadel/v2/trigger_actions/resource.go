package trigger_actions

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func GetResource() *schema.Resource {
	return &schema.Resource{
		Description: "Resource representing triggers, when actions get started",
		Schema: map[string]*schema.Schema{
			orgIDVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the organization",
				ForceNew:    true,
			},
			flowTypeVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Type of the flow to which the action triggers belong",
				ForceNew:    true,
			},
			triggerTypeVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Trigger type on when the actions get triggered",
				ForceNew:    true,
			},
			actionsVar: {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Required:    true,
				Description: "IDs of the triggered actions",
			},
		},
		DeleteContext: delete,
		CreateContext: create,
		UpdateContext: update,
		ReadContext:   read,
		Importer:      &schema.ResourceImporter{StateContext: schema.ImportStatePassthroughContext},
	}
}
