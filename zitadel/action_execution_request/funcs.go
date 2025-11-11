package action_execution_request

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/action/v2"

	actionexecutionbase "github.com/zitadel/terraform-provider-zitadel/v2/zitadel/action_execution_base"
)

func buildCondition(d *schema.ResourceData) (*action.Condition, error) {
	req := &action.RequestExecution{}
	if method, ok := d.GetOk(MethodVar); ok {
		req.Condition = &action.RequestExecution_Method{Method: method.(string)}
	} else if service, ok := d.GetOk(ServiceVar); ok {
		req.Condition = &action.RequestExecution_Service{Service: service.(string)}
	} else if all, ok := d.GetOk(AllVar); ok && all.(bool) {
		req.Condition = &action.RequestExecution_All{All: true}
	} else {
		return nil, fmt.Errorf("invalid request condition: must set one of method, service, or all")
	}
	return &action.Condition{ConditionType: &action.Condition_Request{Request: req}}, nil
}

func IdFromConditionFn(condition *action.Condition) (string, error) {
	computeID := func(value string) string {
		if value == "" {
			return "request"
		}
		if strings.HasPrefix(value, "/") {
			return "request" + value
		}
		return "request/" + value
	}

	if req := condition.GetRequest(); req != nil {
		if method := req.GetMethod(); method != "" {
			return computeID(method), nil
		} else if service := req.GetService(); service != "" {
			return computeID(service), nil
		} else if req.GetAll() {
			return computeID(""), nil
		}
	}
	return "", fmt.Errorf("unknown condition type for ID generation: %v", condition.GetConditionType())
}

func readExecution(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	execution, diags := actionexecutionbase.ReadExecutionBase(ctx, d, m, IdFromConditionFn)
	if diags != nil || execution == nil {
		return diags
	}

	req := execution.GetCondition().GetRequest()
	if req == nil {
		d.SetId("")
		return nil
	}

	if method := req.GetMethod(); method != "" {
		if err := d.Set(MethodVar, method); err != nil {
			return diag.FromErr(err)
		}
	} else if service := req.GetService(); service != "" {
		if err := d.Set(ServiceVar, service); err != nil {
			return diag.FromErr(err)
		}
	} else if req.GetAll() {
		if err := d.Set(AllVar, true); err != nil {
			return diag.FromErr(err)
		}
	}
	if err := d.Set(actionexecutionbase.TargetIDsVar, execution.GetTargets()); err != nil {
		return diag.FromErr(err)
	}
	return nil
}
