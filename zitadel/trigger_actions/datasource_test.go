package trigger_actions_test

import (
	"fmt"
	"testing"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper/test_utils"
)

func TestAccTriggerActionsDatasource(t *testing.T) {
	frame := test_utils.NewOrgTestFrame(t, "zitadel_trigger_actions")
	config := fmt.Sprintf(`
data "zitadel_trigger_actions" "default" {
  org_id       = data.zitadel_org.default.id
  flow_type    = "FLOW_TYPE_EXTERNAL_AUTHENTICATION"
  trigger_type = "TRIGGER_TYPE_POST_AUTHENTICATION"
}`)

	test_utils.RunDatasourceTest(
		t,
		frame.BaseTestFrame,
		config,
		[]string{frame.AsOrgDefaultDependency},
		nil,
		map[string]string{},
	)
}
