package action_execution_base

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/action/v2"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
)

func ReadExecutionBase(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
	idFromCondition IdFromConditionFunc,
) (*action.Execution, diag.Diagnostics) {
	tflog.Info(ctx, "started read")

	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return nil, diag.Errorf("failed to get client")
	}

	client, err := helper.GetActionClient(ctx, clientinfo)
	if err != nil {
		return nil, diag.FromErr(err)
	}

	id := helper.GetID(d, IDVar)
	var found *action.Execution
	// An execution set moments ago may not be listed by the query API yet, so wait for it to appear.
	_, err = helper.RetryUntilFound(ctx, func() (bool, error) {
		resp, err := client.ListExecutions(ctx, &action.ListExecutionsRequest{})
		if err != nil {
			return false, fmt.Errorf("failed to list executions: %w", err)
		}
		for _, execution := range resp.GetExecutions() {
			idPtr, err := idFromCondition(execution.GetCondition())
			if err != nil {
				return false, err
			}
			if idPtr == nil {
				// different execution type → skip
				continue
			}
			if *idPtr == id {
				found = execution
				return true, nil
			}
		}
		return false, nil
	})
	if err != nil {
		return nil, diag.FromErr(err)
	}

	if found == nil || len(found.GetTargets()) == 0 {
		d.SetId("")
		return nil, nil
	}
	d.SetId(id)
	return found, nil
}
