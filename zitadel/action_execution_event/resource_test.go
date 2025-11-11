package action_execution_event_test

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/action/v2"

	actionexecutionbase "github.com/zitadel/terraform-provider-zitadel/v2/zitadel/action_execution_base"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/action_target"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper/test_utils"
)

// TestAccActionExecution_Event asserts the lifecycle of an event
// execution. It checks creation, update, and import.
func TestAccActionExecution_Event(t *testing.T) {
	frame := test_utils.NewInstanceTestFrame(t, "zitadel_action_execution_event")
	targetFrame := test_utils.NewInstanceTestFrame(t, "zitadel_action_target")
	targetFrame2 := test_utils.NewInstanceTestFrame(t, "zitadel_action_target_2")

	targetResource := createTargetResource(t, targetFrame, "default")
	targetResource2 := createTargetResource(t, targetFrame2, "default_2")

	eventName := "user.human.added"
	executionID := "event/" + eventName
	importID := "event:" + eventName
	executionIDRegex := regexp.MustCompile(fmt.Sprintf(`^%s$`, regexp.QuoteMeta(executionID)))

	baseHCL := `
resource "zitadel_action_execution_event" "default" {
  event      = "%s"
  target_ids = %s
}
`
	exampleProperty := `[zitadel_action_target.default.id]`
	updatedProperty := `[zitadel_action_target.default.id, zitadel_action_target.default_2.id]`

	resourceFunc := func(targetIDsProperty string, secret string) string {
		if targetIDsProperty == updatedProperty {
			return fmt.Sprintf(`
%s
%s
%s
`, targetResource, targetResource2, fmt.Sprintf(baseHCL, eventName, targetIDsProperty))
		}
		return fmt.Sprintf(`
%s
%s
`, targetResource, fmt.Sprintf(baseHCL, eventName, targetIDsProperty))
	}

	test_utils.RunLifecyleTest(
		t,
		frame.BaseTestFrame,
		[]string{},
		resourceFunc,
		exampleProperty,
		updatedProperty,
		"", "", "",
		true,
		checkRemoteExecution(frame, executionID),
		executionIDRegex,
		test_utils.CheckIsNotFoundFromPropertyCheck(checkRemoteExecution(frame, executionID), "execution not found"),
		test_utils.ChainImportStateIdFuncs(
			func(state *terraform.State) (string, error) {
				return importID, nil
			},
		),
		"",
	)
}

func createTargetResource(t *testing.T, frame *test_utils.InstanceTestFrame, resourceName string) string {
	targetResource, targetAttrs := test_utils.ReadExample(t, test_utils.Resources, "zitadel_action_target")
	targetResource = strings.Replace(targetResource, `"default"`, `"`+resourceName+`"`, 1)
	nameAttribute := test_utils.AttributeValue(t, action_target.NameVar, targetAttrs).AsString()
	return strings.Replace(targetResource, nameAttribute, frame.UniqueResourcesID, 1)
}

func checkRemoteExecution(frame *test_utils.InstanceTestFrame, expectedID string) func(string) resource.TestCheckFunc {
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
				currentID, err := actionexecutionbase.IdFromCondition(execution.GetCondition())
				if err != nil {
					return err
				}
				if currentID == expectedID {
					return nil
				}
			}

			return fmt.Errorf("execution not found")
		}
	}
}
