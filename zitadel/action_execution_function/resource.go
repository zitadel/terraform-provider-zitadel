package action_execution_function

import (
	"context"
	"fmt"

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
		CreateContext: actionexecutionbase.NewSetExecution(buildCondition),
		DeleteContext: actionexecutionbase.NewDeleteExecution(buildCondition),
		ReadContext:   readExecution,
		UpdateContext: actionexecutionbase.NewSetExecution(buildCondition),
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				functionName := d.Id()
				internalID := fmt.Sprintf("function/%s", functionName)
				d.SetId(internalID)

				if err := d.Set(NameVar, functionName); err != nil {
					return nil, err
				}
				return []*schema.ResourceData{d}, nil
			},
		},
	}
}
