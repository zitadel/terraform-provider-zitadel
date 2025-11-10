package action_execution_test

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/action_execution"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/action_target/action_target_test_dep"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper/test_utils"
	actionv2 "github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/action/v2"
)

func TestAccExecution(t *testing.T) {
	os.Setenv("TF_ACC", "1")
	frame := test_utils.NewOrgTestFrame(t, "zitadel_action_execution")
	actionDep, _ := action_target_test_dep.Create(t, frame)

	resourceExample, exampleAttributes := test_utils.ReadExample(t, test_utils.Resources, frame.ResourceType)

	exampleProperty := test_utils.AttributeValue(t, action_execution.EventGroupVar, exampleAttributes).AsString()
	updatedProperty := "action"

	test_utils.RunLifecyleTest(
		t,
		frame.BaseTestFrame,
		[]string{frame.AsOrgDefaultDependency, actionDep},
		test_utils.ReplaceAll(resourceExample, exampleProperty, ""),
		exampleProperty, updatedProperty,
		"", "", "",
		true,
		checkRemoteProperty(frame),
		regexp.MustCompile(fmt.Sprintf(`^%s$`, ".+")),
		checkDestroy(frame),
		nil,
	)
}

func checkRemoteProperty(frame *test_utils.OrgTestFrame) func(string) resource.TestCheckFunc {
	return func(expect string) resource.TestCheckFunc {
		return func(state *terraform.State) error {
			client, err := helper.GetActionClient(context.Background(), frame.ClientInfo)
			if err != nil {
				return fmt.Errorf("failed to get client: %w", err)
			}
			remoteExecutions, err := client.ListExecutions(
				context.Background(),
				&actionv2.ListExecutionsRequest{},
			)
			if err != nil {
				return err
			}
			if len(remoteExecutions.GetExecutions()) == 0 {
				return fmt.Errorf("no executions found")
			}
			actual := remoteExecutions.GetExecutions()[0].Condition.GetEvent().GetGroup()
			if expect != actual {
				return fmt.Errorf("expected event %s, but got %s", expect, actual)
			}
			return nil
		}
	}
}

func checkDestroy(frame *test_utils.OrgTestFrame) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		client, err := helper.GetActionClient(context.Background(), frame.ClientInfo)
		if err != nil {
			return fmt.Errorf("failed to get client: %w", err)
		}
		remoteExecutions, err := client.ListExecutions(
			context.Background(),
			&actionv2.ListExecutionsRequest{},
		)
		if err != nil {
			return err
		}
		if len(remoteExecutions.GetExecutions()) != 0 {
			return fmt.Errorf("expected no executions, but found %d", len(remoteExecutions.GetExecutions()))
		}
		return nil
	}
}
