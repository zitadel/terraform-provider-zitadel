package action_execution_function

import (
	"context"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	actionexecutionbase "github.com/zitadel/terraform-provider-zitadel/v2/zitadel/action_execution_base"
)

func GetResource() *schema.Resource {
	return &schema.Resource{
		Description: "Resource representing an action execution triggered by a function.",
		Schema: map[string]*schema.Schema{
			NameVar: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the function, e.g. `Action.Flow.Type.ExternalAuthentication.Action.TriggerType.PostAuthentication`",
			},
			actionexecutionbase.TargetIDsVar: {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				Description: "The list of target IDs to call.",
			},
		},
		CreateContext: actionexecutionbase.NewSetExecution(buildCondition, IdFromConditionFn),
		DeleteContext: actionexecutionbase.NewDeleteExecution(buildCondition),
		ReadContext:   readExecution,
		UpdateContext: actionexecutionbase.NewSetExecution(buildCondition, IdFromConditionFn),
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
				m := regexp.MustCompile(`^([^:\s]+)$`).FindStringSubmatch(d.Id())
				if m == nil {
					return nil, fmt.Errorf("invalid import ID: %s. Must be a function name like 'Action.Flow.Type.ExternalAuthentication.Action.TriggerType.PostAuthentication'", d.Id())
				}
				d.SetId("function/" + m[1])
				return []*schema.ResourceData{d}, nil
			},
		},
	}
}
