package action_execution_function

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/action/v2"

	actionexecutionbase "github.com/zitadel/terraform-provider-zitadel/v2/zitadel/action_execution_base"
)

func buildCondition(d *schema.ResourceData) (*action.Condition, error) {
	condition := &action.Condition{
		ConditionType: &action.Condition_Function{
			Function: &action.FunctionExecution{
				Name: d.Get(NameVar).(string),
			},
		},
	}
	return condition, nil
}

func IdFromConditionFn(condition *action.Condition) (string, error) {
	if fn := condition.GetFunction(); fn != nil {
		if name := fn.GetName(); name != "" {
			return "function/" + name, nil
		}
		return "function", nil
	}
	return "", fmt.Errorf("unknown condition type for ID generation: %v", condition.GetConditionType())
}

func readExecution(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	execution, diags := actionexecutionbase.ReadExecutionBase(ctx, d, m, IdFromConditionFn)
	if diags != nil || execution == nil {
		return diags
	}

	fn := execution.GetCondition().GetFunction()
	if fn == nil {
		d.SetId("")
		return nil
	}

	if err := d.Set(NameVar, fn.GetName()); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set(actionexecutionbase.TargetIDsVar, execution.GetTargets()); err != nil {
		return diag.FromErr(err)
	}
	return nil
}
