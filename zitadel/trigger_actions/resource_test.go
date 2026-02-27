package trigger_actions_test

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/management"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/action/action_test_dep"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper/test_utils"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/trigger_actions"
)

func TestAccTriggerActions(t *testing.T) {
	frame := test_utils.NewOrgTestFrame(t, "zitadel_trigger_actions")
	resourceExample, exampleAttributes := test_utils.ReadExample(t, test_utils.Resources, frame.ResourceType)
	exampleProperty := test_utils.AttributeValue(t, trigger_actions.TriggerTypeVar, exampleAttributes).AsString()
	flowType := test_utils.AttributeValue(t, trigger_actions.FlowTypeVar, exampleAttributes).AsString()
	actionDep, _ := action_test_dep.Create(t, frame)
	updatedProperty := "TRIGGER_TYPE_PRE_USERINFO_CREATION"
	test_utils.RunLifecyleTest(
		t,
		frame.BaseTestFrame,
		[]string{frame.AsOrgDefaultDependency, actionDep},
		test_utils.ReplaceAll(resourceExample, exampleProperty, ""),
		exampleProperty, updatedProperty,
		"", "", "",
		false,
		checkRemoteProperty(frame, flowType),
		regexp.MustCompile(fmt.Sprintf("^%s_([A-Z_]+)_(%s|%s)$", helper.ZitadelGeneratedIdPattern, exampleProperty, updatedProperty)),
		test_utils.CheckIsNotFoundFromPropertyCheck(checkRemoteProperty(frame, flowType), exampleProperty),
		test_utils.ChainImportStateIdFuncs(
			test_utils.ImportStateAttribute(frame.BaseTestFrame, trigger_actions.FlowTypeVar),
			test_utils.ImportStateAttribute(frame.BaseTestFrame, trigger_actions.TriggerTypeVar),
			test_utils.ImportOrgId(frame),
		),
	)
}

func TestAccTriggerActionsExternalAuthFlow(t *testing.T) {
	frame := test_utils.NewOrgTestFrame(t, "zitadel_trigger_actions")
	actionDep, actionID := action_test_dep.Create(t, frame)

	resourceConfig := fmt.Sprintf(`
%s
%s
resource "zitadel_trigger_actions" "default" {
  org_id       = data.zitadel_org.default.id
  flow_type    = "FLOW_TYPE_EXTERNAL_AUTHENTICATION"
  trigger_type = "TRIGGER_TYPE_POST_AUTHENTICATION"
  action_ids   = ["%s"]
}
`, frame.ProviderSnippet, frame.AsOrgDefaultDependency, actionID)

	_ = actionDep

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: frame.V6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: resourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(frame.TerraformName, "flow_type", "FLOW_TYPE_EXTERNAL_AUTHENTICATION"),
					resource.TestCheckResourceAttr(frame.TerraformName, "trigger_type", "TRIGGER_TYPE_POST_AUTHENTICATION"),
					checkRemoteProperty(frame, "FLOW_TYPE_EXTERNAL_AUTHENTICATION")("TRIGGER_TYPE_POST_AUTHENTICATION"),
				),
			},
		},
	})
}

func checkRemoteProperty(frame *test_utils.OrgTestFrame, flowType string) func(string) resource.TestCheckFunc {
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
