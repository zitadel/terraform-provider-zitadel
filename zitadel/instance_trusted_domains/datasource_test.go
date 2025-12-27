package instance_trusted_domains_test

import (
	"testing"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper/test_utils"
)

func TestAccInstanceTrustedDomainsDatasource(t *testing.T) {
	frame := test_utils.NewInstanceTestFrame(t, "zitadel_instance_trusted_domains")
	test_utils.RunDatasourceTest(
		t,
		frame.BaseTestFrame,
		`
		data "zitadel_instance_trusted_domains" "default" {}
		`,
		nil,
		nil,
		map[string]string{},
	)
}
