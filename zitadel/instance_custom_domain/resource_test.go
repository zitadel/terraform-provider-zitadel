package instance_custom_domain_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper/test_utils"
)

func TestAccInstanceCustomDomain(t *testing.T) {
	t.Skip("Skipping test - system API user in Zitadel v4.11.0 lacks system.domain.write permission (AUTH-5mWD2)")
	systemFrame := test_utils.NewSystemTestFrame(t, "zitadel_instance_custom_domain")
	instanceFrame := test_utils.NewInstanceTestFrame(t, "zitadel_instance_custom_domain")

	resourceConfig := func(domain string, _ string) string {
		return fmt.Sprintf(`
resource "zitadel_instance_custom_domain" "default" {
    instance_id = "%s"
    domain      = "%s"
}
`, instanceFrame.InstanceID, domain)
	}

	test_utils.RunLifecyleTest(
		t,
		*systemFrame,
		nil,
		resourceConfig,
		"login.example.com",
		"login.example2.com",
		"",
		"",
		"",
		false,
		func(expect string) resource.TestCheckFunc {
			return test_utils.CheckNothing
		},
		regexp.MustCompile(`^.+$`),
		test_utils.CheckNothing,
		test_utils.ChainImportStateIdFuncs(
			test_utils.ImportResourceId(*systemFrame),
		),
	)
}
