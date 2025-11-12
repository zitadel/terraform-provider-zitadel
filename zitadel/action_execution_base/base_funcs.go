package action_execution

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/action/v2"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
)

const (
	TargetIDsVar = "target_ids"
)

func ReadExecutionBase(ctx context.Context, d *schema.ResourceData, m interface{}, idFromCondition IdFromConditionFunc) (*action.Execution, diag.Diagnostics) {
	tflog.Info(ctx, "started read")

	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return nil, diag.Errorf("failed to get client")
	}

	client, err := helper.GetActionClient(ctx, clientinfo)
	if err != nil {
		return nil, diag.FromErr(err)
	}

	resp, err := client.ListExecutions(ctx, &action.ListExecutionsRequest{})
	if err != nil {
		return nil, diag.Errorf("failed to list executions: %v", err)
	}

	for _, execution := range resp.GetExecutions() {
		currentID, err := idFromCondition(execution.GetCondition())
		if err != nil {
			// Expected: execution is a different type, skip it
			// Only specific type-matching errors are expected
			continue
		}

		if currentID == d.Id() {
			if len(execution.GetTargets()) == 0 {
				d.SetId("")
				return nil, nil
			}
			return execution, nil
		}
	}

	d.SetId("")
	return nil, nil
}

func ConditionFromID(id string) (*action.Condition, error) {
	parts := strings.SplitN(id, "/", 2)
	if len(parts) == 0 {
		return nil, fmt.Errorf("invalid execution ID format: %s", id)
	}

	conditionType := parts[0]
	value := ""
	if len(parts) == 2 {
		value = parts[1]
	}

	if conditionType == "request" && strings.HasPrefix(id, "request/") {
		value = "/" + value
	}
	if conditionType == "response" && strings.HasPrefix(id, "response/") {
		value = "/" + value
	}

	condition := &action.Condition{}
	switch conditionType {
	case "request":
		req := &action.RequestExecution{}
		if value == "" {
			req.Condition = &action.RequestExecution_All{All: true}
		} else if strings.HasPrefix(value, "/") {
			req.Condition = &action.RequestExecution_Method{Method: value}
		} else {
			req.Condition = &action.RequestExecution_Service{Service: value}
		}
		condition.ConditionType = &action.Condition_Request{Request: req}
	case "response":
		resp := &action.ResponseExecution{}
		if value == "" {
			resp.Condition = &action.ResponseExecution_All{All: true}
		} else if strings.HasPrefix(value, "/") {
			resp.Condition = &action.ResponseExecution_Method{Method: value}
		} else {
			resp.Condition = &action.ResponseExecution_Service{Service: value}
		}
		condition.ConditionType = &action.Condition_Response{Response: resp}
	case "function":
		condition.ConditionType = &action.Condition_Function{Function: &action.FunctionExecution{Name: value}}
	case "event":
		event := &action.EventExecution{}
		if value == "" {
			event.Condition = &action.EventExecution_All{All: true}
		} else if strings.HasSuffix(value, ".*") {
			event.Condition = &action.EventExecution_Group{Group: value}
		} else {
			event.Condition = &action.EventExecution_Event{Event: value}
		}
		condition.ConditionType = &action.Condition_Event{Event: event}
	default:
		return nil, fmt.Errorf("unknown condition type in ID: %s", conditionType)
	}
	return condition, nil
}
