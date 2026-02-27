package action_execution_function_test

import (
	"fmt"
	"testing"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper/test_utils"
)

func TestAccActionExecutionFunctionDatasource(t *testing.T) {
	frame := test_utils.NewInstanceTestFrame(t, "zitadel_action_execution_function")
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
resource "zitadel_action_execution_function" "default" {
  name       = "preaccesstoken"
  target_ids = [zitadel_action_target.default.id]
}`

	config := `
data "zitadel_action_execution_function" "default" {
  id = zitadel_action_execution_function.default.id
}`

	test_utils.RunDatasourceTest(
		t,
		frame.BaseTestFrame,
		config,
		[]string{targetDep, executionDep},
		nil,
		map[string]string{
			"name": "preaccesstoken",
		},
	)
}
