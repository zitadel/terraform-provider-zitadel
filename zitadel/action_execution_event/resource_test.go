package action_execution_event_test

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/action/v2"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/action_execution_base/test_helpers"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/action_execution_event"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/action_target"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper/test_utils"
)

func TestSchemaConsistency(t *testing.T) {
	resource := action_execution_event.GetResource()
	if resource == nil {
		t.Fatal("GetResource() returned nil")
	}

	const isWritable = true

	t.Run("ExactlyOneOf", func(t *testing.T) {
		test_helpers.CheckSchemaExactlyOneOfConsistency(t, resource.Schema)
	})

	t.Run("InternalValidate", func(t *testing.T) {
		test_helpers.CheckSchemaInternalValidation(t, resource, isWritable)
	})
}

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
		test_helpers.CheckRemoteExecution(frame, executionID, deriveEventID),
		executionIDRegex,
		test_utils.CheckIsNotFoundFromPropertyCheck(test_helpers.CheckRemoteExecution(frame, executionID, deriveEventID), ""),
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

func deriveEventID(cond *action.Condition) (string, error) {
	ev := cond.GetEvent()
	if ev == nil {
		return "", fmt.Errorf("no event condition")
	}
	if e := ev.GetEvent(); e != "" {
		return "event/" + e, nil
	}
	if g := ev.GetGroup(); g != "" {
		if !strings.HasSuffix(g, ".*") {
			g += ".*"
		}
		return "event/" + g, nil
	}
	if ev.GetAll() {
		return "event", nil
	}
	return "", fmt.Errorf("unknown event condition")
}
