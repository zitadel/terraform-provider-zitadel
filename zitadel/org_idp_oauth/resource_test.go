package org_idp_oauth_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper/test_utils"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/idp_utils"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/org_idp_utils/org_idp_test_utils"
)

func TestAccOrgIdPOAuth(t *testing.T) {
	org_idp_test_utils.RunOrgLifecyleTest(t, "zitadel_org_idp_oauth", idp_utils.ClientSecretVar)
}

func TestAccOrgIdPOAuthUsePKCE(t *testing.T) {
	frame := test_utils.NewOrgTestFrame(t, "zitadel_org_idp_oauth")
	resourceConfig := fmt.Sprintf(`
%s
%s
resource "zitadel_org_idp_oauth" "default" {
  org_id                 = data.zitadel_org.default.id
  name                   = "%s"
  client_id              = "test_client_id"
  client_secret          = "test_client_secret"
  authorization_endpoint = "https://example.com/oauth/authorize"
  token_endpoint         = "https://example.com/oauth/token"
  user_endpoint          = "https://example.com/oauth/userinfo"
  id_attribute           = "user_id"
  scopes                 = ["openid", "profile"]
  is_linking_allowed     = false
  is_creation_allowed    = true
  is_auto_creation       = false
  is_auto_update         = true
  use_pkce               = true
}`, frame.ProviderSnippet, frame.AsOrgDefaultDependency, frame.UniqueResourcesID)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: frame.V6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: resourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(frame.TerraformName, "use_pkce", "true"),
				),
			},
		},
	})
}
