package system_features_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper/test_utils"
)

func TestAccSystemFeatures(t *testing.T) {
	frame := test_utils.NewSystemTestFrame(t, "zitadel_system_features")

	resourceExample := `
resource "zitadel_system_features" "default" {
	login_default_org = true
}
	`

	resourceExampleUpdated := `
resource "zitadel_system_features" "default" {
	login_default_org = false
	user_schema = true
}
	`

	test_utils.RunLifecyleTest(
		t,
		*frame,
		nil,
		func(property, secret string) string {
			if property == resourceExample {
				return resourceExample
			}
			return resourceExampleUpdated
		},
		resourceExample,
		resourceExampleUpdated,
		"", "", "",
		false,
		func(_ string) resource.TestCheckFunc {
			return func(_ *terraform.State) error { return nil }
		},
		regexp.MustCompile("^system$"),
		func(_ *terraform.State) error { return nil },
		test_utils.ChainImportStateIdFuncs(
			func(_ *terraform.State) (string, error) { return "system", nil },
		),
	)
}
