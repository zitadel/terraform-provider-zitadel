package action_execution_event

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	actionexecutionbase "github.com/zitadel/terraform-provider-zitadel/v2/zitadel/action_execution_base"
)

func GetDatasource() *schema.Resource {
	return actionexecutionbase.NewActionExecutionDatasource(
		"Datasource representing an action execution triggered by an event.",
		"The ID of this resource. Must be set to the condition, e.g. `event:user.human.added`, `group:user.human`, or `all`/`all:`",
		map[string]*schema.Schema{
			EventVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The specific event to trigger on.",
			},
			GroupVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The event group to trigger on.",
			},
			AllVar: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Trigger on all events.",
			},
		},
		readExecution,
	)
}
