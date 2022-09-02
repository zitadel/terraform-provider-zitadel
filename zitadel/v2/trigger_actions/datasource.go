package trigger_actions

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func GetDatasource() *schema.Resource {
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
			},
			triggerTypeVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Trigger type on when the actions get triggered",
			},
			actionsVar: {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed:    true,
				Description: "IDs of the triggered actions",
			},
		},
		ReadContext: read,
		Importer:    &schema.ResourceImporter{StateContext: schema.ImportStatePassthroughContext},
	}
}
