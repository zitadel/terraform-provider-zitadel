package action_execution_response_test

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/action/v2"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/action_execution_base/test_helpers"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/action_target"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper/test_utils"
)

func TestAccActionExecution_Response(t *testing.T) {
	frame := test_utils.NewInstanceTestFrame(t, "zitadel_action_execution_response")
	targetFrame := test_utils.NewInstanceTestFrame(t, "zitadel_action_target")
	targetFrame2 := test_utils.NewInstanceTestFrame(t, "zitadel_action_target_2")

	targetResource := createTargetResource(t, targetFrame, "default")
	targetResource2 := createTargetResource(t, targetFrame2, "default_2")

	method := "/zitadel.session.v2.SessionService/GetSession"
	executionID := "response" + method
	importID := "method:" + method
	executionIDRegex := regexp.MustCompile(fmt.Sprintf(`^%s$`, regexp.QuoteMeta(executionID)))

	baseHCL := `
resource "zitadel_action_execution_response" "default" {
  method       = "%s"
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
`, targetResource, targetResource2, fmt.Sprintf(baseHCL, method, targetIDsProperty))
		}
		return fmt.Sprintf(`
%s
%s
`, targetResource, fmt.Sprintf(baseHCL, method, targetIDsProperty))
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
		test_helpers.CheckRemoteExecution(frame, executionID, deriveResponseID),
		executionIDRegex,
		test_utils.CheckIsNotFoundFromPropertyCheck(test_helpers.CheckRemoteExecution(frame, executionID, deriveResponseID), ""),
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

func deriveResponseID(cond *action.Condition) (string, error) {
	resp := cond.GetResponse()
	if resp == nil {
		return "", fmt.Errorf("no response condition")
	}
	if m := resp.GetMethod(); m != "" {
		return "response" + m, nil
	}
	if s := resp.GetService(); s != "" {
		return "response/" + s, nil
	}
	if resp.GetAll() {
		return "response", nil
	}
	return "", fmt.Errorf("unknown response condition")
}
