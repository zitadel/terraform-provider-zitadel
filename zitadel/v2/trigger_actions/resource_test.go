package trigger_actions_test

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/management"
	"google.golang.org/protobuf/types/known/durationpb"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/helper"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/helper/test_utils"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/trigger_actions"
)

func TestAccTriggerActions(t *testing.T) {
	resourceName := "zitadel_trigger_actions"
	flowType := "FLOW_TYPE_CUSTOMISE_TOKEN"
	initialProperty := "TRIGGER_TYPE_PRE_ACCESS_TOKEN_CREATION"
	updatedProperty := "TRIGGER_TYPE_PRE_USERINFO_CREATION"
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
	test_utils.RunLifecyleTest[string](
		t,
		frame.BaseTestFrame,
		func(configProperty, _ string) string {
			return fmt.Sprintf(`
resource "%s" "%s" {
	org_id              = "%s"
flow_type = "%s"
  trigger_type = "%s"
  action_ids   = ["%s"]
}`, resourceName, frame.UniqueResourcesID, frame.OrgID, flowType, configProperty, action.GetId())
		},
		initialProperty, updatedProperty,
		"", "", "",
		false,
		checkRemoteProperty(*frame, flowType),
		regexp.MustCompile(fmt.Sprintf("^%s_([A-Z_]+)_(%s|%s)$", helper.ZitadelGeneratedIdPattern, initialProperty, updatedProperty)),
		test_utils.CheckIsNotFoundFromPropertyCheck(checkRemoteProperty(*frame, flowType), initialProperty),
		nil,
	)
}

func checkRemoteProperty(frame test_utils.OrgTestFrame, flowType string) func(string) resource.TestCheckFunc {
	return func(expect string) resource.TestCheckFunc {
		return func(state *terraform.State) error {
			flowTypeValues := helper.EnumValueMap(trigger_actions.FlowTypes())
			resp, err := frame.GetFlow(frame, &management.GetFlowRequest{Type: strconv.Itoa(int(flowTypeValues[flowType]))})
			if err != nil {
				return fmt.Errorf("flow type not found: %w", err)
			}
			typesMapping := trigger_actions.TriggerTypes()
			var foundTypes []string
			for _, actual := range resp.GetFlow().GetTriggerActions() {
				idInt, err := strconv.Atoi(actual.GetTriggerType().GetId())
				if err != nil {
					return err
				}
				foundType := typesMapping[int32(idInt)]
				foundTypes = append(foundTypes, foundType)
				if foundType == expect {
					return nil
				}
			}
			return fmt.Errorf("expected trigger type %s not found in %v: %w", expect, foundTypes, test_utils.ErrNotFound)
		}
	}
}
