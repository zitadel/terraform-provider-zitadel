package default_oidc_settings_test

import (
	"testing"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper/test_utils"
)

func TestAccDefaultOidcSettingsDatasource(t *testing.T) {
	frame := test_utils.NewInstanceTestFrame(t, "zitadel_default_oidc_settings")
	test_utils.RunDatasourceTest(
		t,
		frame.BaseTestFrame,
		`data "zitadel_default_oidc_settings" "default" {}`,
		nil,
		nil,
		map[string]string{},
	)
}
