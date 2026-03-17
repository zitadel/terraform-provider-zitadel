package instance_secret_generator_test

import (
	"testing"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper/test_utils"
)

func TestAccInstanceSecretGeneratorDatasource(t *testing.T) {
	frame := test_utils.NewInstanceTestFrame(t, "zitadel_instance_secret_generator")
	test_utils.RunDatasourceTest(
		t,
		frame.BaseTestFrame,
		`data "zitadel_instance_secret_generator" "default" {
  generator_type = "verify_email_code"
}`,
		nil,
		nil,
		map[string]string{},
	)
}
