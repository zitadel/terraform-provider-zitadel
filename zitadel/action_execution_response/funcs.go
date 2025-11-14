package action_execution_response

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
	resp := &action.ResponseExecution{}
	if method, ok := d.GetOk(MethodVar); ok {
		resp.Condition = &action.ResponseExecution_Method{Method: method.(string)}
	} else if service, ok := d.GetOk(ServiceVar); ok {
		resp.Condition = &action.ResponseExecution_Service{Service: service.(string)}
	} else if all, ok := d.GetOk(AllVar); ok && all.(bool) {
		resp.Condition = &action.ResponseExecution_All{All: true}
	} else {
		return nil, fmt.Errorf("invalid response condition: must set one of method, service, or all")
	}
	return &action.Condition{ConditionType: &action.Condition_Response{Response: resp}}, nil
}

func IdFromConditionFn(condition *action.Condition) (*string, error) {
	computeID := func(value string) string {
		if value == "" { // all responses
			return "response"
		}
		if strings.HasPrefix(value, "/") { // method-based condition
			return "response" + value
		}
		return "response/" + value // service-based condition
	}

	if resp := condition.GetResponse(); resp == nil { // not a response execution → skip
		return nil, nil
	} else if method := resp.GetMethod(); method != "" { // method-based condition
		id := computeID(method)
		return &id, nil
	} else if service := resp.GetService(); service != "" { // service-based condition
		id := computeID(service)
		return &id, nil
	} else if resp.GetAll() { // all responses
		id := computeID("")
		return &id, nil
	} else { // malformed response condition
		return nil, fmt.Errorf("invalid response condition: %#v", resp)
	}
}

func readExecution(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	execution, diags := actionexecutionbase.ReadExecutionBase(ctx, d, m, IdFromConditionFn)
	if diags != nil || execution == nil {
		return diags
	}

	resp := execution.GetCondition().GetResponse()
	if resp == nil {
		d.SetId("")
		return nil
	}

	if method := resp.GetMethod(); method != "" {
		if err := d.Set(MethodVar, method); err != nil {
			return diag.FromErr(err)
		}
	} else if service := resp.GetService(); service != "" {
		if err := d.Set(ServiceVar, service); err != nil {
			return diag.FromErr(err)
		}
	} else if resp.GetAll() {
		if err := d.Set(AllVar, true); err != nil {
			return diag.FromErr(err)
		}
	}
	if err := d.Set(actionexecutionbase.TargetIDsVar, execution.GetTargets()); err != nil {
		return diag.FromErr(err)
	}
	return nil
}
