package application_oidc_client_secret_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/ephemeral/ephtest"
)

// TestAccApplicationOIDCClientSecret verifies the ephemeral resource
// regenerates and returns the client secret of a zitadel_application_oidc app.
func TestAccApplicationOIDCClientSecret(t *testing.T) {
	frame := ephtest.NewOrgFrame(t)
	config := fmt.Sprintf(`
%s
%s
resource "zitadel_project" "default" {
  name   = "proj-%s"
  org_id = data.zitadel_org.default.id
}
resource "zitadel_application_oidc" "default" {
  project_id                = zitadel_project.default.id
  org_id                    = data.zitadel_org.default.id
  name                      = "appoidc-%s"
  redirect_uris             = ["https://localhost.com"]
  response_types            = ["OIDC_RESPONSE_TYPE_CODE"]
  grant_types               = ["OIDC_GRANT_TYPE_AUTHORIZATION_CODE"]
  post_logout_redirect_uris = ["https://localhost.com"]
  app_type                  = "OIDC_APP_TYPE_WEB"
  auth_method_type          = "OIDC_AUTH_METHOD_TYPE_BASIC"
  version                   = "OIDC_VERSION_1_0"
  dev_mode                  = true
}
ephemeral "zitadel_application_oidc_client_secret" "test" {
  project_id = zitadel_project.default.id
  app_id     = zitadel_application_oidc.default.id
  org_id     = data.zitadel_org.default.id
}
provider "echo" {
  data = ephemeral.zitadel_application_oidc_client_secret.test
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
