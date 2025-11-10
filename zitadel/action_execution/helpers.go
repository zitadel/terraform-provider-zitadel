package action_execution

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	actionv2 "github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/action/v2"
)

func buildEventConditionFromResourceData(d *schema.ResourceData) (*actionv2.Condition, error) {
	event := d.Get(EventNameVar).(string)
	group := d.Get(EventGroupVar).(string)
	all := d.Get(AllVar).(bool)
	var eventCond actionv2.EventExecution
	if event != "" {
		eventCond.Condition = &actionv2.EventExecution_Event{Event: event}
	} else if group != "" {
		eventCond.Condition = &actionv2.EventExecution_Group{Group: group}
	} else if all {
		eventCond.Condition = &actionv2.EventExecution_All{All: true}
	} else {
		return nil, fmt.Errorf("for event type, one of event, group, or all must be set")
	}
	return &actionv2.Condition{
		ConditionType: &actionv2.Condition_Event{Event: &eventCond},
	}, nil
}

func buildRequestResponseConditionFromResourceData(d *schema.ResourceData) (*actionv2.Condition, error) {
	service := d.Get(ServiceVar).(string)
	method := d.Get(MethodVar).(string)
	all := d.Get(AllVar).(bool)
	var reqCond actionv2.RequestExecution
	if method != "" {
		reqCond.Condition = &actionv2.RequestExecution_Method{Method: method}
	} else if service != "" {
		reqCond.Condition = &actionv2.RequestExecution_Service{Service: service}
	} else if all {
		reqCond.Condition = &actionv2.RequestExecution_All{All: true}
	} else {
		return nil, fmt.Errorf("for request and response types, one of method, service, or all must be set")
	}
	return &actionv2.Condition{
		ConditionType: &actionv2.Condition_Request{Request: &reqCond},
	}, nil
}

func buildFunctionConditionFromResourceData(d *schema.ResourceData) (*actionv2.Condition, error) {
	function := d.Get(FunctionNameVar).(string)
	if function == "" {
		return nil, fmt.Errorf("for function type, function must be set")
	}
	return &actionv2.Condition{
		ConditionType: &actionv2.Condition_Function{
			Function: &actionv2.FunctionExecution{Name: function},
		},
	}, nil
}

// buildConditionFromResourceData extracts the mutually exclusive condition block from d and returns the appropriate *actionv2.Condition.
func buildConditionFromResourceData(d *schema.ResourceData) (*actionv2.Condition, error) {
	execType, ok := d.GetOk(ExecutionTypeVar)
	if !ok {
		return nil, fmt.Errorf("execution_type must be set")
	}
	switch execType {
	case executionTypeEvent:
		return buildEventConditionFromResourceData(d)
	case executionTypeRequest:
		return buildRequestResponseConditionFromResourceData(d)
	case executionTypeResponse:
		return buildRequestResponseConditionFromResourceData(d)
	case executionTypeFunction:
		return buildFunctionConditionFromResourceData(d)
	default:
		return nil, fmt.Errorf("invalid execution_type: %v", execType)
	}
}
