package trigger_actions

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
)

func GetDatasource() *schema.Resource {
	return &schema.Resource{
		Description: "Resource representing triggers, when actions get started",
		Schema: map[string]*schema.Schema{
			helper.OrgIDVar: helper.OrgIDResourceField,
			FlowTypeVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Type of the flow to which the action triggers belong",
			},
			TriggerTypeVar: {
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
	}
}
