package instance_restrictions_test

import (
	"testing"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper/test_utils"
)

func TestAccInstanceRestrictionsDatasource(t *testing.T) {
	frame := test_utils.NewInstanceTestFrame(t, "zitadel_instance_restrictions")
	test_utils.RunDatasourceTest(
		t,
		frame.BaseTestFrame,
		`data "zitadel_instance_restrictions" "default" {}`,
		nil,
		nil,
		map[string]string{},
	)
}
