package action_execution_event

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/action/v2"

	actionexecutionbase "github.com/zitadel/terraform-provider-zitadel/v2/zitadel/action_execution_base"
)

func buildCondition(d *schema.ResourceData) (*action.Condition, error) {
	event := &action.EventExecution{}
	if eventName, ok := d.GetOk(EventVar); ok {
		event.Condition = &action.EventExecution_Event{Event: eventName.(string)}
	} else if group, ok := d.GetOk(GroupVar); ok {
		event.Condition = &action.EventExecution_Group{Group: group.(string)}
	} else if all, ok := d.GetOk(AllVar); ok && all.(bool) {
		event.Condition = &action.EventExecution_All{All: true}
	} else {
		return nil, fmt.Errorf("invalid event condition: must set one of event, group, or all")
	}
	return &action.Condition{ConditionType: &action.Condition_Event{Event: event}}, nil
}

func IdFromConditionFn(condition *action.Condition) (*string, error) {
	computeID := func(value string) string {
		if value == "" { // all events
			return "event"
		}
		return "event/" + value // event/group-style value
	}

	if event := condition.GetEvent(); event == nil { // not an event execution → skip
		return nil, nil
	} else if eventName := event.GetEvent(); eventName != "" { // specific event
		id := computeID(eventName)
		return &id, nil
	} else if group := event.GetGroup(); group != "" { // event group
		id := computeID(group)
		return &id, nil
	} else if event.GetAll() { // all events
		id := computeID("")
		return &id, nil
	} else { // malformed event condition
		return nil, fmt.Errorf("invalid event condition: %#v", event)
	}
}

func readExecution(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	execution, diags := actionexecutionbase.ReadExecutionBase(ctx, d, m, IdFromConditionFn)
	if diags != nil || execution == nil {
		return diags
	}

	event := execution.GetCondition().GetEvent()
	if event == nil {
		d.SetId("")
		return nil
	}

	if eventName := event.GetEvent(); eventName != "" {
		if err := d.Set(EventVar, eventName); err != nil {
			return diag.FromErr(err)
		}
	} else if group := event.GetGroup(); group != "" {
		if err := d.Set(GroupVar, group); err != nil {
			return diag.FromErr(err)
		}
	} else if event.GetAll() {
		if err := d.Set(AllVar, true); err != nil {
			return diag.FromErr(err)
		}
	}
	if err := d.Set(actionexecutionbase.TargetIDsVar, execution.GetTargets()); err != nil {
		return diag.FromErr(err)
	}
	return nil
}
