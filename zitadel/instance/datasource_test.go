package instance_test

import (
	"fmt"
	"testing"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper/test_utils"
)

func TestAccInstanceDatasource(t *testing.T) {
	frame := test_utils.NewInstanceTestFrame(t, "zitadel_instance")
	test_utils.RunDatasourceTest(
		t,
		frame.BaseTestFrame,
		fmt.Sprintf(`
		data "zitadel_instance" "default" {
			instance_id = "%s"
		}
		`, frame.InstanceID),
		nil,
		nil,
		map[string]string{
			"name": "instance-level-tests",
		},
	)
}
