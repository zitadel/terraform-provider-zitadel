package action_execution_event_test

import (
	"fmt"
	"testing"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper/test_utils"
)

func TestAccActionExecutionEventDatasource_Event(t *testing.T) {
	frame := test_utils.NewInstanceTestFrame(t, "zitadel_action_execution_event")
	targetDep := fmt.Sprintf(`
resource "zitadel_action_target" "default" {
  name               = "%s"
  endpoint           = "https://example.com"
  target_type        = "REST_ASYNC"
  timeout            = "10s"
  interrupt_on_error = false
  payload_type       = "PAYLOAD_TYPE_JSON"
}`, frame.UniqueResourcesID)

	executionDep := `
resource "zitadel_action_execution_event" "default" {
  event      = "user.human.added"
  target_ids = [zitadel_action_target.default.id]
}`

	config := `
data "zitadel_action_execution_event" "default" {
  id = zitadel_action_execution_event.default.id
}`

	test_utils.RunDatasourceTest(
		t,
		frame.BaseTestFrame,
		config,
		[]string{targetDep, executionDep},
		nil,
		map[string]string{
			"event": "user.human.added",
		},
	)
}

func TestAccActionExecutionEventDatasource_Group(t *testing.T) {
	frame := test_utils.NewInstanceTestFrame(t, "zitadel_action_execution_event")
	targetDep := fmt.Sprintf(`
resource "zitadel_action_target" "default" {
  name               = "%s"
  endpoint           = "https://example.com"
  target_type        = "REST_ASYNC"
  timeout            = "10s"
  interrupt_on_error = false
  payload_type       = "PAYLOAD_TYPE_JSON"
}`, frame.UniqueResourcesID)

	executionDep := `
resource "zitadel_action_execution_event" "default" {
  group      = "user.human"
  target_ids = [zitadel_action_target.default.id]
}`

	config := `
data "zitadel_action_execution_event" "default" {
  id = zitadel_action_execution_event.default.id
}`

	test_utils.RunDatasourceTest(
		t,
		frame.BaseTestFrame,
		config,
		[]string{targetDep, executionDep},
		nil,
		map[string]string{
			"group": "user.human",
		},
	)
}

func TestAccActionExecutionEventDatasource_All(t *testing.T) {
	frame := test_utils.NewInstanceTestFrame(t, "zitadel_action_execution_event")
	targetDep := fmt.Sprintf(`
resource "zitadel_action_target" "default" {
  name               = "%s"
  endpoint           = "https://example.com"
  target_type        = "REST_ASYNC"
  timeout            = "10s"
  interrupt_on_error = false
  payload_type       = "PAYLOAD_TYPE_JSON"
}`, frame.UniqueResourcesID)

	executionDep := `
resource "zitadel_action_execution_event" "default" {
  all        = true
  target_ids = [zitadel_action_target.default.id]
}`

	config := `
data "zitadel_action_execution_event" "default" {
  id = zitadel_action_execution_event.default.id
}`

	test_utils.RunDatasourceTest(
		t,
		frame.BaseTestFrame,
		config,
		[]string{targetDep, executionDep},
		nil,
		map[string]string{
			"all": "true",
		},
	)
}
