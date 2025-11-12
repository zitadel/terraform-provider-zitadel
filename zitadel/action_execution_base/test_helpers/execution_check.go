package test_helpers

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/action/v2"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper/test_utils"
)

type IdFromConditionFunc func(condition *action.Condition) (string, error)

func CheckRemoteExecution(frame *test_utils.InstanceTestFrame, expectedID string, idFromCondition IdFromConditionFunc) func(string) resource.TestCheckFunc {
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
				currentID, err := idFromCondition(execution.GetCondition())
				if err != nil {
					continue // Skip executions of other types
				}
				if currentID == expectedID {
					return nil
				}
			}

			return test_utils.ErrNotFound
		}
	}
}
