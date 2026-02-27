package idp_github_es_test

import (
	"fmt"
	"testing"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper/test_utils"
)

func TestAccIdpGithubEsDatasource(t *testing.T) {
	frame := test_utils.NewInstanceTestFrame(t, "zitadel_idp_github_es")
	resourceDep := fmt.Sprintf(`
resource "zitadel_idp_github_es" "default" {
  name                   = "%s"
  client_id              = "dummy"
  client_secret          = "dummy"
  scopes                 = ["openid", "profile", "email"]
  authorization_endpoint = "https://github.example.com/login/oauth/authorize"
  token_endpoint         = "https://github.example.com/login/oauth/access_token"
  user_endpoint          = "https://api.github.example.com/user"
  is_linking_allowed     = false
  is_creation_allowed    = true
  is_auto_creation       = false
  is_auto_update         = true
}`, frame.UniqueResourcesID)

	config := `
data "zitadel_idp_github_es" "default" {
  id = zitadel_idp_github_es.default.id
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
