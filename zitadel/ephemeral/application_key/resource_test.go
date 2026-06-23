package application_key_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/ephemeral/ephtest"
)

// TestAccApplicationKey verifies the ephemeral resource creates a new JSON key
// for an API application and returns the key material. The type name
// zitadel_application_key is shared with the managed resource; this also
// confirms the muxed provider serves both forms without conflict.
func TestAccApplicationKey(t *testing.T) {
	frame := ephtest.NewOrgFrame(t)
	config := fmt.Sprintf(`
%s
%s
resource "zitadel_project" "default" {
  name   = "proj-%s"
  org_id = data.zitadel_org.default.id
}
resource "zitadel_application_api" "default" {
  org_id           = data.zitadel_org.default.id
  project_id       = zitadel_project.default.id
  name             = "appapi-%s"
  auth_method_type = "API_AUTH_METHOD_TYPE_PRIVATE_KEY_JWT"
}
ephemeral "zitadel_application_key" "test" {
  project_id      = zitadel_project.default.id
  app_id          = zitadel_application_api.default.id
  org_id          = data.zitadel_org.default.id
  expiration_date = "2519-04-01T08:45:00Z"
}
provider "echo" {
  data = ephemeral.zitadel_application_key.test
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
