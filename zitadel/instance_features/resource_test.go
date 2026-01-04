// zitadel/instance_features/resource_test.go
package instance_features_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper/test_utils"
)

func TestAccInstanceFeatures(t *testing.T) {
	frame := test_utils.NewInstanceTestFrame(t, "zitadel_instance_features")

	resourceExample := `
resource "zitadel_instance_features" "default" {
	login_default_org = true
}
	`

	test_utils.RunLifecyleTest(
		t,
		frame.BaseTestFrame,
		nil,
		test_utils.ReplaceAll(resourceExample, "", ""),
		"", "",
		"", "", "",
		false,
		func(_ string) resource.TestCheckFunc {
			return func(_ *terraform.State) error { return nil }
		},
		nil,
		func(_ *terraform.State) error { return nil },
		test_utils.ChainImportStateIdFuncs(
			func(_ *terraform.State) (string, error) { return "instance", nil },
		),
	)
}
