package trigger_actions_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/management"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/action/action_test_dep"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/helper"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/helper/test_utils"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/trigger_actions"
)

func TestAccTriggerActions(t *testing.T) {
	frame := test_utils.NewOrgTestFrame(t, "zitadel_trigger_actions")
	resourceExample, exampleAttributes := test_utils.ReadExample(t, test_utils.Resources, frame.ResourceType)
	exampleProperty := test_utils.AttributeValue(t, trigger_actions.TriggerTypeVar, exampleAttributes).AsString()
	flowType := test_utils.AttributeValue(t, trigger_actions.FlowTypeVar, exampleAttributes).AsString()
	actionDep, _ := action_test_dep.Create(t, frame)
	test_utils.RunLifecyleTest(
		t,
		frame.BaseTestFrame,
		[]string{frame.AsOrgDefaultDependency, actionDep},
		test_utils.ReplaceAll(resourceExample, exampleProperty, ""),
		exampleProperty, "TRIGGER_TYPE_PRE_USERINFO_CREATION",
		"", "",
		false,
		checkRemoteProperty(*frame, flowType),
		test_utils.ZITADEL_GENERATED_ID_REGEX,
		test_utils.CheckIsNotFoundFromPropertyCheck(checkRemoteProperty(*frame, flowType), exampleProperty),
		nil, nil, "", "",
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
