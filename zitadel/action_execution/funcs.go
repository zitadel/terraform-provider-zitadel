package action_execution

import (
	"context"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	actionv2 "github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/action/v2"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
)

func delete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started delete")

	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := helper.GetActionClient(ctx, clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}

	cond, err := buildConditionFromResourceData(d)
	if err != nil {
		return diag.FromErr(err)
	}
	// To delete an execution we need to perform a set execution with no targets
	_, err = client.SetExecution(ctx, &actionv2.SetExecutionRequest{
		Condition: cond,
		Targets:   []string{},
	})
	if err != nil {
		return diag.Errorf("failed to delete execution: %v", err)
	}
	return nil
}

func update(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started update")

	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := helper.GetActionClient(ctx, clientinfo)
	if err != nil {
		return diag.Errorf("failed to get action client: %v", err)
	}
	cond, err := buildConditionFromResourceData(d)
	if err != nil {
		return diag.Errorf("failed to build condition from resource data: %v", err)
	}

	targetsInterface := d.Get(TargetsVar).([]interface{})
	targets := make([]string, len(targetsInterface))
	for i, t := range targetsInterface {
		targets[i] = t.(string)
	}

	_, err = client.SetExecution(ctx, &actionv2.SetExecutionRequest{
		Condition: cond,
		Targets:   targets,
	})
	if err != nil {
		return diag.Errorf("failed to update execution: %v", err)
	}

	return nil
}

func create(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started create")

	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := helper.GetActionClient(ctx, clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}
	cond, err := buildConditionFromResourceData(d)
	if err != nil {
		return diag.FromErr(err)
	}

	targetsInterface := d.Get(TargetsVar).([]interface{})
	targets := make([]string, len(targetsInterface))
	for i, t := range targetsInterface {
		targets[i] = t.(string)
	}

	req := &actionv2.SetExecutionRequest{
		Condition: cond,
		Targets:   targets,
	}
	_, err = client.SetExecution(ctx, req)
	if err != nil {
		return diag.Errorf("failed to create execution: %v", err)
	}

	// "event:{event:\"user.human.added\"}"
	d.SetId(cond.String())
	return nil
}

func read(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started read")

	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := helper.GetActionClient(ctx, clientinfo)
	if err != nil {
		return diag.Errorf("failed to get action client: %v", err)
	}

	cond, err := buildConditionFromResourceData(d)
	if err != nil {
		return diag.Errorf("failed to build condition from resource data: %v", err)
	}

	foundExec, err := findExecutionByCondition(ctx, client, cond)
	if err != nil {
		return diag.Errorf("failed to find execution by condition: %v", err)
	}

	if foundExec == nil {
		tflog.Info(ctx, "execution not found, removing from state")
		d.SetId("")
		return nil
	}

	set := map[string]interface{}{
		TargetsVar: foundExec.GetTargets(),
	}

	if foundExec.Condition.GetEvent() != nil {
		set[ExecutionTypeVar] = executionTypeEvent
		set[EventNameVar] = foundExec.Condition.GetEvent().GetEvent()
		set[EventGroupVar] = foundExec.Condition.GetEvent().GetGroup()
		set[AllVar] = foundExec.Condition.GetEvent().GetAll()
	} else if foundExec.Condition.GetFunction() != nil {
		set[ExecutionTypeVar] = executionTypeFunction
		set[FunctionNameVar] = foundExec.Condition.GetFunction().GetName()
	} else if foundExec.Condition.GetRequest() != nil {
		set[ExecutionTypeVar] = executionTypeRequest
		set[ServiceVar] = foundExec.Condition.GetRequest().GetService()
		set[MethodVar] = foundExec.Condition.GetRequest().GetMethod()
		set[AllVar] = foundExec.Condition.GetRequest().GetAll()
	} else if foundExec.Condition.GetResponse() != nil {
		set[ExecutionTypeVar] = executionTypeResponse
		set[ServiceVar] = foundExec.Condition.GetResponse().GetService()
		set[MethodVar] = foundExec.Condition.GetResponse().GetMethod()
		set[AllVar] = foundExec.Condition.GetResponse().GetAll()
	}

	for k, v := range set {
		if err := d.Set(k, v); err != nil {
			return diag.Errorf("failed to set %s of execution: %v", k, err)
		}
	}

	d.SetId(cond.String())

	return nil
}

func findExecutionByCondition(ctx context.Context, client actionv2.ActionServiceClient, cond *actionv2.Condition) (*actionv2.Execution, error) {
	execResp, err := client.ListExecutions(ctx, &actionv2.ListExecutionsRequest{})
	if err != nil {
		return nil, err
	}

	var foundExec *actionv2.Execution
	for _, exec := range execResp.Executions {
		if exec.GetCondition().String() == cond.String() {
			foundExec = exec
			break
		}
	}

	return foundExec, nil
}
