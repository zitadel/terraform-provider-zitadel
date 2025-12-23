package zitadel_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper/test_utils"
)

func TestAccSessionTokenDatasource(t *testing.T) {
	frame := test_utils.NewOrgTestFrame(t, "zitadel")
	config := `data "zitadel" "token" {}`
	test_utils.RunDatasourceTest(
		t,
		frame.BaseTestFrame,
		config,
		nil,
		resource.TestCheckResourceAttrSet("data.zitadel.token", "token"),
		nil,
	)
}
