package action_test

import (
	"fmt"
	"testing"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper/test_utils"
)

func TestAccActionDatasource(t *testing.T) {
	frame := test_utils.NewOrgTestFrame(t, "zitadel_action")
	resourceDep := fmt.Sprintf(`
resource "zitadel_action" "default" {
  org_id          = data.zitadel_org.default.id
  name            = "%s"
  script          = "function(ctx, api) { /* noop */ }"
  timeout         = "10s"
  allowed_to_fail = true
}`, frame.UniqueResourcesID)

	config := fmt.Sprintf(`
data "zitadel_action" "default" {
  org_id    = data.zitadel_org.default.id
  action_id = zitadel_action.default.id
}`)

	test_utils.RunDatasourceTest(
		t,
		frame.BaseTestFrame,
		config,
		[]string{frame.AsOrgDefaultDependency, resourceDep},
		nil,
		map[string]string{
			"name": frame.UniqueResourcesID,
		},
	)
}
