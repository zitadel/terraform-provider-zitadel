package trigger_actions_test

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"google.golang.org/protobuf/types/known/durationpb"

	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/management"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/helper/test_utils"
)

func TestAccTriggerActions(t *testing.T) {
	resourceName := "zitadel_trigger_actions"
	flowType := "FLOW_TYPE_CUSTOMISE_TOKEN"
	initialTriggerType := "TRIGGER_TYPE_PRE_ACCESS_TOKEN_CREATION"
	updatedTriggerType := "TRIGGER_TYPE_POST_AUTHENTICATION"
	frame, err := test_utils.NewOrgTestFrame(resourceName)
	if err != nil {
		t.Fatalf("setting up test context failed: %v", err)
	}
	// Always creates a new action
	action, err := frame.CreateAction(frame, &management.CreateActionRequest{
		Name:          frame.UniqueResourcesID,
		Script:        "not a script",
		Timeout:       durationpb.New(10 * time.Second),
		AllowedToFail: true,
	})
	if err != nil {
		t.Fatalf("failed to create action: %v", err)
	}
	test_utils.RunLifecyleTest(
		t,
		frame.BaseTestFrame,
		func(name, _ string) string {
			return fmt.Sprintf(`
resource "%s" "%s" {
	org_id              = "%s"
flow_type = "%s"
  trigger_type = "%s"
  action_ids   = ["%s"]
}`, resourceName, frame.UniqueResourcesID, frame.OrgID, flowType, name, action.GetId())
		},
		initialTriggerType, updatedTriggerType,
		"", "",
		CheckTriggerType(*frame, flowType),
		CheckDestroy(*frame, flowType, []string{initialTriggerType, updatedTriggerType}),
		nil, nil, "", "",
	)
}

var errTriggerTypeNotFound = errors.New("trigger type not found")

func CheckTriggerType(frame test_utils.OrgTestFrame, flowType string) func(string) resource.TestCheckFunc {
	return func(expectTriggerType string) resource.TestCheckFunc {
		return func(state *terraform.State) error {
			triggerTypes, err := frame.ListFlowTriggerTypes(frame, &management.ListFlowTriggerTypesRequest{Type: flowType})
			if err != nil {
				return err
			}
			result := triggerTypes.GetResult()
			for _, actual := range result {
				if actual.GetId() == expectTriggerType {
					return nil
				}
			}
			return fmt.Errorf("expected trigger type %s not found in %v: %w", expectTriggerType, result, errTriggerTypeNotFound)
		}
	}
}

func CheckDestroy(frame test_utils.OrgTestFrame, flowType string, testTypes []string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		for _, testTriggerType := range testTypes {
			if err := CheckTriggerType(frame, flowType)(testTriggerType)(state); !errors.Is(err, errTriggerTypeNotFound) {
				return fmt.Errorf("expected error %v, but got %v", errTriggerTypeNotFound, err)
			}
		}
		return nil
	}
}
