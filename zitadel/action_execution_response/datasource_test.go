package action_execution_response_test

import (
	"fmt"
	"testing"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper/test_utils"
)

func TestAccActionExecutionResponseDatasource_All(t *testing.T) {
	frame := test_utils.NewInstanceTestFrame(t, "zitadel_action_execution_response")
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
resource "zitadel_action_execution_response" "default" {
  all        = true
  target_ids = [zitadel_action_target.default.id]
}`

	config := `
data "zitadel_action_execution_response" "default" {
  id = zitadel_action_execution_response.default.id
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

func TestAccActionExecutionResponseDatasource_Method(t *testing.T) {
	frame := test_utils.NewInstanceTestFrame(t, "zitadel_action_execution_response")
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
resource "zitadel_action_execution_response" "default" {
  method     = "/zitadel.session.v2.SessionService/ListSessions"
  target_ids = [zitadel_action_target.default.id]
}`

	config := `
data "zitadel_action_execution_response" "default" {
  id = zitadel_action_execution_response.default.id
}`

	test_utils.RunDatasourceTest(
		t,
		frame.BaseTestFrame,
		config,
		[]string{targetDep, executionDep},
		nil,
		map[string]string{
			"method": "/zitadel.session.v2.SessionService/ListSessions",
		},
	)
}

func TestAccActionExecutionResponseDatasource_Service(t *testing.T) {
	frame := test_utils.NewInstanceTestFrame(t, "zitadel_action_execution_response")
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
resource "zitadel_action_execution_response" "default" {
  service    = "zitadel.session.v2.SessionService"
  target_ids = [zitadel_action_target.default.id]
}`

	config := `
data "zitadel_action_execution_response" "default" {
  id = zitadel_action_execution_response.default.id
}`

	test_utils.RunDatasourceTest(
		t,
		frame.BaseTestFrame,
		config,
		[]string{targetDep, executionDep},
		nil,
		map[string]string{
			"service": "zitadel.session.v2.SessionService",
		},
	)
}
