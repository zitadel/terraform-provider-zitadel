package action_execution_request_test

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

func TestAccActionExecution_Request(t *testing.T) {
	frame := test_utils.NewInstanceTestFrame(t, "zitadel_action_execution_request")
	targetFrame := test_utils.NewInstanceTestFrame(t, "zitadel_action_target")
	targetFrame2 := test_utils.NewInstanceTestFrame(t, "zitadel_action_target_2")

	targetResource := createTargetResource(t, targetFrame, "default")
	targetResource2 := createTargetResource(t, targetFrame2, "default_2")

	method := "/zitadel.session.v2.SessionService/ListSessions"
	executionID := "request" + method
	importID := "method:" + method
	executionIDRegex := regexp.MustCompile(fmt.Sprintf(`^%s$`, regexp.QuoteMeta(executionID)))

	baseHCL := `
resource "zitadel_action_execution_request" "default" {
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
		test_helpers.CheckRemoteExecution(frame, executionID, deriveRequestID),
		executionIDRegex,
		test_utils.CheckIsNotFoundFromPropertyCheck(test_helpers.CheckRemoteExecution(frame, executionID, deriveRequestID), ""),
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

func deriveRequestID(cond *action.Condition) (string, error) {
	req := cond.GetRequest()
	if req == nil {
		return "", fmt.Errorf("no request condition")
	}
	if m := req.GetMethod(); m != "" {
		return "request" + m, nil
	}
	if s := req.GetService(); s != "" {
		return "request/" + s, nil
	}
	if req.GetAll() {
		return "request", nil
	}
	return "", fmt.Errorf("unknown request condition")
}
