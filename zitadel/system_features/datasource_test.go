package system_features_test

import (
	"testing"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper/test_utils"
)

func TestAccSystemFeaturesDatasource(t *testing.T) {
	frame := test_utils.NewSystemTestFrame(t, "zitadel_system_features")
	test_utils.RunDatasourceTest(
		t,
		*frame,
		`data "zitadel_system_features" "default" {}`,
		nil,
		nil,
		map[string]string{},
	)
}
