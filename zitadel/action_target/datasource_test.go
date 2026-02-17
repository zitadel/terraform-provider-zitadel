package action_target_test

import (
	"fmt"
	"testing"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper/test_utils"
)

func TestAccActionTargetDatasource(t *testing.T) {
	t.Skip("Skipping: flaky due to eventual consistency after action_execution tests on same instance")
	frame := test_utils.NewInstanceTestFrame(t, "zitadel_action_target")
	resourceDep := fmt.Sprintf(`
resource "zitadel_action_target" "default" {
  name               = "%s"
  endpoint           = "https://example.com/datasource-test"
  target_type        = "REST_WEBHOOK"
  timeout            = "10s"
  interrupt_on_error = false
  payload_type       = "PAYLOAD_TYPE_JSON"
}
`, frame.UniqueResourcesID)

	config := `
data "zitadel_action_target" "default" {
  target_id = zitadel_action_target.default.id
}
`

	test_utils.RunDatasourceTest(
		t,
		frame.BaseTestFrame,
		config,
		[]string{resourceDep},
		nil,
		map[string]string{
			"name":               frame.UniqueResourcesID,
			"endpoint":           "https://example.com/datasource-test",
			"target_type":        "REST_WEBHOOK",
			"timeout":            "10s",
			"interrupt_on_error": "false",
			"payload_type":       "PAYLOAD_TYPE_JSON",
		},
	)
}

func TestAccActionTargetDatasourcePayloadTypeJWT(t *testing.T) {
	frame := test_utils.NewInstanceTestFrame(t, "zitadel_action_target")
	resourceDep := fmt.Sprintf(`
resource "zitadel_action_target" "default" {
  name               = "%s"
  endpoint           = "https://example.com/datasource-test"
  target_type        = "REST_ASYNC"
  timeout            = "5s"
  interrupt_on_error = false
  payload_type       = "PAYLOAD_TYPE_JWT"
}
`, frame.UniqueResourcesID)

	config := `
data "zitadel_action_target" "default" {
  target_id = zitadel_action_target.default.id
}
`

	test_utils.RunDatasourceTest(
		t,
		frame.BaseTestFrame,
		config,
		[]string{resourceDep},
		nil,
		map[string]string{
			"payload_type": "PAYLOAD_TYPE_JWT",
		},
	)
}

func TestAccActionTargetDatasourcePayloadTypeJWE(t *testing.T) {
	frame := test_utils.NewInstanceTestFrame(t, "zitadel_action_target")
	resourceDep := fmt.Sprintf(`
resource "zitadel_action_target" "default" {
  name               = "%s"
  endpoint           = "https://example.com/datasource-test"
  target_type        = "REST_ASYNC"
  timeout            = "5s"
  interrupt_on_error = false
  payload_type       = "PAYLOAD_TYPE_JWE"
}
`, frame.UniqueResourcesID)

	config := `
data "zitadel_action_target" "default" {
  target_id = zitadel_action_target.default.id
}
`

	test_utils.RunDatasourceTest(
		t,
		frame.BaseTestFrame,
		config,
		[]string{resourceDep},
		nil,
		map[string]string{
			"payload_type": "PAYLOAD_TYPE_JWE",
		},
	)
}

func TestAccActionTargetDatasourceTargetTypeRestWebhook(t *testing.T) {
	frame := test_utils.NewInstanceTestFrame(t, "zitadel_action_target")
	resourceDep := fmt.Sprintf(`
resource "zitadel_action_target" "default" {
  name               = "%s"
  endpoint           = "https://example.com/datasource-test"
  target_type        = "REST_WEBHOOK"
  timeout            = "10s"
  interrupt_on_error = true
  payload_type       = "PAYLOAD_TYPE_JSON"
}
`, frame.UniqueResourcesID)

	config := `
data "zitadel_action_target" "default" {
  target_id = zitadel_action_target.default.id
}
`

	test_utils.RunDatasourceTest(
		t,
		frame.BaseTestFrame,
		config,
		[]string{resourceDep},
		nil,
		map[string]string{
			"target_type": "REST_WEBHOOK",
		},
	)
}

func TestAccActionTargetDatasourceTargetTypeRestCall(t *testing.T) {
	frame := test_utils.NewInstanceTestFrame(t, "zitadel_action_target")
	resourceDep := fmt.Sprintf(`
resource "zitadel_action_target" "default" {
  name               = "%s"
  endpoint           = "https://example.com/datasource-test"
  target_type        = "REST_CALL"
  timeout            = "10s"
  interrupt_on_error = true
  payload_type       = "PAYLOAD_TYPE_JSON"
}
`, frame.UniqueResourcesID)

	config := `
data "zitadel_action_target" "default" {
  target_id = zitadel_action_target.default.id
}
`

	test_utils.RunDatasourceTest(
		t,
		frame.BaseTestFrame,
		config,
		[]string{resourceDep},
		nil,
		map[string]string{
			"target_type": "REST_CALL",
		},
	)
}

func TestAccActionTargetDatasourceTargetTypeRestAsync(t *testing.T) {
	frame := test_utils.NewInstanceTestFrame(t, "zitadel_action_target")
	resourceDep := fmt.Sprintf(`
resource "zitadel_action_target" "default" {
  name               = "%s"
  endpoint           = "https://example.com/datasource-test"
  target_type        = "REST_ASYNC"
  timeout            = "10s"
  interrupt_on_error = false
  payload_type       = "PAYLOAD_TYPE_JSON"
}
`, frame.UniqueResourcesID)

	config := `
data "zitadel_action_target" "default" {
  target_id = zitadel_action_target.default.id
}
`

	test_utils.RunDatasourceTest(
		t,
		frame.BaseTestFrame,
		config,
		[]string{resourceDep},
		nil,
		map[string]string{
			"target_type": "REST_ASYNC",
		},
	)
}
