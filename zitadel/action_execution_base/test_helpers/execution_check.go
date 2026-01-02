package test_helpers

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/action/v2"

	actionexecutionbase "github.com/zitadel/terraform-provider-zitadel/v2/zitadel/action_execution_base"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper/test_utils"
)

func CheckRemoteExecution(
	frame *test_utils.InstanceTestFrame,
	expectedID string,
	idFromCondition actionexecutionbase.IdFromConditionFunc,
) func(string) resource.TestCheckFunc {
	return func(targetsCount string) resource.TestCheckFunc {
		return func(state *terraform.State) error {
			client, err := helper.GetActionClient(context.Background(), frame.ClientInfo)
			if err != nil {
				return fmt.Errorf("failed to get client: %w", err)
			}

			resp, err := client.ListExecutions(context.Background(), &action.ListExecutionsRequest{})
			if err != nil {
				return fmt.Errorf("failed to list executions: %w", err)
			}

			for _, execution := range resp.GetExecutions() {
				idPtr, err := idFromCondition(execution.GetCondition())
				if err != nil {
					return fmt.Errorf("failed to derive execution id: %w", err)
				}
				if idPtr != nil && *idPtr == expectedID {
					return nil
				}
			}

			return test_utils.ErrNotFound
		}
	}
}
