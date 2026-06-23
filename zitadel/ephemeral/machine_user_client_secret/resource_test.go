package machine_user_client_secret_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/ephemeral/ephtest"
)

// TestAccMachineUserClientSecret verifies the ephemeral resource generates and
// returns a fresh client_id/client_secret pair for a machine user.
func TestAccMachineUserClientSecret(t *testing.T) {
	frame := ephtest.NewOrgFrame(t)
	config := fmt.Sprintf(`
%s
%s
resource "zitadel_machine_user" "default" {
  org_id      = data.zitadel_org.default.id
  user_name   = "machine-%s@example.com"
  name        = "machine-%s"
  description = "ephemeral secret test user"
  with_secret = false
}
ephemeral "zitadel_machine_user_client_secret" "test" {
  user_id = zitadel_machine_user.default.id
  org_id  = data.zitadel_org.default.id
}
provider "echo" {
  data = ephemeral.zitadel_machine_user_client_secret.test
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
					resource.TestCheckResourceAttrSet("echo.test", "data.client_id"),
					resource.TestCheckResourceAttrSet("echo.test", "data.client_secret"),
				),
			},
		},
	})
}
