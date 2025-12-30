package instance_custom_domains_test

import (
	"testing"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper/test_utils"
)

func TestAccInstanceCustomDomainsDatasource(t *testing.T) {
	frame := test_utils.NewInstanceTestFrame(t, "zitadel_instance_custom_domains")
	test_utils.RunDatasourceTest(
		t,
		frame.BaseTestFrame,
		`
		data "zitadel_instance_custom_domains" "default" {}
		`,
		nil,
		nil,
		map[string]string{},
	)
}
