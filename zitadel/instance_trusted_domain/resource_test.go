package instance_trusted_domain_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper/test_utils"
)

func TestAccInstanceTrustedDomain(t *testing.T) {
	frame := test_utils.NewInstanceTestFrame(t, "zitadel_instance_trusted_domain")

	resourceConfig := func(domain string, _ string) string {
		return fmt.Sprintf(`
resource "zitadel_instance_trusted_domain" "default" {
    instance_id = data.zitadel_instance.current.id
    domain      = "%s"
}

data "zitadel_instance" "current" {}
`, domain)
	}

	test_utils.RunLifecyleTest(
		t,
		frame.BaseTestFrame,
		nil,
		resourceConfig,
		"idp.partner.com",
		"idp.partner2.com",
		"",
		"",
		"",
		false,
		func(expect string) resource.TestCheckFunc {
			return test_utils.CheckNothing
		},
		regexp.MustCompile(`^.+$`),
		test_utils.CheckNothing,
		nil,
	)
}
