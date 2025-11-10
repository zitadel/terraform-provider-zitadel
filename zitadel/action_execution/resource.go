package action_execution

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func GetResource() *schema.Resource {
	return &schema.Resource{
		Description: "Resource representing an execution, which defines when and which targets are executed based on triggers.",
		Schema: map[string]*schema.Schema{
			"execution_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The type of execution trigger. One of: events, request, response, function.",
			},
			// Event fields
			"event": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The event name to trigger on (for events type).",
			},
			"group": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The event group (for events type).",
			},
			"all": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to trigger on all events in the group (for events, request and response types).",
			},
			// Request/Response fields
			"service": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The gRPC service name (for request/response type).",
			},
			"method": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The gRPC method name (for request/response type).",
			},
			// Function fields
			"function": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The function name to trigger on (for function type).",
			},
			// Targets
			TargetsVar: {
				Type:        schema.TypeList,
				Required:    true,
				Description: "The ordered list of target IDs to execute.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
		CreateContext: create,
		DeleteContext: delete,
		ReadContext:   read,
		UpdateContext: update,
		// No Importer: executions are identified by condition, not ID
	}
}
