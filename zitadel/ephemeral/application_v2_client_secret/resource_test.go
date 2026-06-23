package application_v2_client_secret_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/ephemeral/ephtest"
)

// TestAccApplicationV2ClientSecret verifies that the ephemeral resource
// regenerates and returns the client secret of a zitadel_application_v2 OIDC
// application. The echo provider captures the ephemeral value so we can assert
// it was produced; the secret itself is never written to state.
func TestAccApplicationV2ClientSecret(t *testing.T) {
	frame := ephtest.NewOrgFrame(t)
	config := fmt.Sprintf(`
%s
%s
resource "zitadel_project" "default" {
  name   = "proj-%s"
  org_id = data.zitadel_org.default.id
}
resource "zitadel_application_v2" "default" {
  project_id = zitadel_project.default.id
  org_id     = data.zitadel_org.default.id
  name       = "appv2-%s"
  oidc {
    redirect_uris                = ["https://localhost.com"]
    response_types               = ["OIDC_RESPONSE_TYPE_CODE"]
    grant_types                  = ["OIDC_GRANT_TYPE_AUTHORIZATION_CODE"]
    post_logout_redirect_uris    = ["https://localhost.com"]
    app_type                     = "OIDC_APP_TYPE_WEB"
    auth_method_type             = "OIDC_AUTH_METHOD_TYPE_BASIC"
    version                      = "OIDC_VERSION_1_0"
    clock_skew                   = "0s"
    dev_mode                     = true
    access_token_type            = "OIDC_TOKEN_TYPE_BEARER"
    access_token_role_assertion  = false
    id_token_role_assertion      = false
    id_token_userinfo_assertion  = false
    additional_origins           = []
    skip_native_app_success_page = false
  }
}
ephemeral "zitadel_application_v2_client_secret" "test" {
  project_id     = zitadel_project.default.id
  application_id = zitadel_application_v2.default.id
  org_id         = data.zitadel_org.default.id
}
provider "echo" {
  data = ephemeral.zitadel_application_v2_client_secret.test
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
