package machine_key_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/ephemeral/ephtest"
)

// TestAccMachineKey verifies the ephemeral resource creates a new JSON key for a
// machine user and returns the key material. Note the ephemeral resource shares
// the type name zitadel_machine_key with the managed resource; this test also
// confirms the muxed provider serves both an `ephemeral` and a `resource` of
// that name without conflict.
func TestAccMachineKey(t *testing.T) {
	frame := ephtest.NewOrgFrame(t)
	config := fmt.Sprintf(`
%s
%s
resource "zitadel_machine_user" "default" {
  org_id      = data.zitadel_org.default.id
  user_name   = "machine-%s@example.com"
  name        = "machine-%s"
  description = "ephemeral key test user"
  with_secret = false
}
ephemeral "zitadel_machine_key" "test" {
  user_id         = zitadel_machine_user.default.id
  org_id          = data.zitadel_org.default.id
  expiration_date = "2519-04-01T08:45:00Z"
}
provider "echo" {
  data = ephemeral.zitadel_machine_key.test
}
resource "echo" "test" {}
`, frame.ProviderSnippet, frame.OrgDependency, frame.UniqueID, frame.UniqueID)

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_10_0),
		},
		ProtoV6ProviderFactories: frame.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("echo.test", "data.key_id"),
					resource.TestCheckResourceAttrSet("echo.test", "data.key_details"),
				),
			},
		},
	})
}
