package application_api_client_secret_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/ephemeral/ephtest"
)

// TestAccApplicationAPIClientSecret verifies the ephemeral resource regenerates
// and returns the client secret of a zitadel_application_api app.
func TestAccApplicationAPIClientSecret(t *testing.T) {
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
  auth_method_type = "API_AUTH_METHOD_TYPE_BASIC"
}
ephemeral "zitadel_application_api_client_secret" "test" {
  project_id = zitadel_project.default.id
  app_id     = zitadel_application_api.default.id
  org_id     = data.zitadel_org.default.id
}
provider "echo" {
  data = ephemeral.zitadel_application_api_client_secret.test
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
					resource.TestCheckResourceAttrSet("echo.test", "data.client_secret"),
				),
			},
		},
	})
}
