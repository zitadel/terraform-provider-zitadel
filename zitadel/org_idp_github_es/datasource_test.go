package org_idp_github_es_test

import (
	"fmt"
	"testing"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper/test_utils"
)

func TestAccOrgIdpGithubEsDatasource(t *testing.T) {
	frame := test_utils.NewOrgTestFrame(t, "zitadel_org_idp_github_es")
	resourceDep := fmt.Sprintf(`
resource "zitadel_org_idp_github_es" "default" {
  org_id                 = data.zitadel_org.default.id
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
data "zitadel_org_idp_github_es" "default" {
  org_id = data.zitadel_org.default.id
  id     = zitadel_org_idp_github_es.default.id
}`

	test_utils.RunDatasourceTest(
		t,
		frame.BaseTestFrame,
		config,
		[]string{frame.AsOrgDefaultDependency, resourceDep},
		nil,
		map[string]string{
			"name": frame.UniqueResourcesID,
		},
	)
}
