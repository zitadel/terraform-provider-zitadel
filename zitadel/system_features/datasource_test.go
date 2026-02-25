package system_features_test

import (
	"testing"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper/test_utils"
)

func TestAccSystemFeaturesDatasource(t *testing.T) {
	frame := test_utils.NewInstanceTestFrame(t, "zitadel_system_features")
	test_utils.RunDatasourceTest(
		t,
		frame.BaseTestFrame,
		`data "zitadel_system_features" "default" {}`,
		nil,
		nil,
		map[string]string{},
	)
}
