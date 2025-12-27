package instance_test

import (
	"testing"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper/test_utils"
)

func TestAccInstanceDatasource(t *testing.T) {
	frame := test_utils.NewInstanceTestFrame(t, "zitadel_instance")
	test_utils.RunDatasourceTest(
		t,
		frame.BaseTestFrame,
		`
		data "zitadel_instance" "default" {}
		`,
		nil,
		nil,
		map[string]string{
			"name": "instance-level-tests",
		},
	)
}
