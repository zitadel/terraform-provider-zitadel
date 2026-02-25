package idp_oauth_test

import (
	"fmt"
	"testing"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper/test_utils"
)

func TestAccIdpOauthDatasource(t *testing.T) {
	frame := test_utils.NewInstanceTestFrame(t, "zitadel_idp_oauth")
	resourceDep := fmt.Sprintf(`
resource "zitadel_idp_oauth" "default" {
  name                   = "%s"
  client_id              = "dummy"
  client_secret          = "dummy"
  scopes                 = ["openid", "profile", "email"]
  authorization_endpoint = "https://oauth.example.com/authorize"
  token_endpoint         = "https://oauth.example.com/token"
  user_endpoint          = "https://oauth.example.com/userinfo"
  id_attribute           = "sub"
  is_linking_allowed     = false
  is_creation_allowed    = true
  is_auto_creation       = false
  is_auto_update         = true
}`, frame.UniqueResourcesID)

	config := `
data "zitadel_idp_oauth" "default" {
  id = zitadel_idp_oauth.default.id
}`

	test_utils.RunDatasourceTest(
		t,
		frame.BaseTestFrame,
		config,
		[]string{resourceDep},
		nil,
		map[string]string{
			"name": frame.UniqueResourcesID,
		},
	)
}
