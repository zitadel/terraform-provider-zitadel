package action_execution_response

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	actionexecutionbase "github.com/zitadel/terraform-provider-zitadel/v2/zitadel/action_execution_base"
)

func GetDatasource() *schema.Resource {
	return actionexecutionbase.NewActionExecutionDatasource(
		"Datasource representing an action execution triggered by a response.",
		"The ID of this resource. Must be set to the condition, e.g. `method:/zitadel.session.v2.SessionService/ListSessions`, `service:zitadel.session.v2.SessionService`, or `all`/`all:`",
		map[string]*schema.Schema{
			MethodVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The gRPC method to trigger on.",
			},
			ServiceVar: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The gRPC service to trigger on.",
			},
			AllVar: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Trigger on all responses.",
			},
		},
		readExecution,
	)
}
