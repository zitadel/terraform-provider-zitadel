package action_execution

import (
	"context"

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
