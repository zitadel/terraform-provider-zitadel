package idp_test

import (
	"fmt"
	"testing"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper/test_utils"
)

func TestAccIdpsDatasource_All(t *testing.T) {
	frame := test_utils.NewInstanceTestFrame(t, "zitadel_idps")
	googleDep, githubDep := idpDeps(frame)

	config := fmt.Sprintf(`
data "zitadel_idps" "default" {
  name        = "%s"
  name_method = "TEXT_QUERY_METHOD_CONTAINS"
  depends_on  = [zitadel_idp_google.default, zitadel_idp_github.default]
}`, frame.UniqueResourcesID)

	test_utils.RunDatasourceTest(
		t,
		frame.BaseTestFrame,
		config,
		[]string{googleDep, githubDep},
		nil,
		map[string]string{
			"ids.#": "2",
		},
	)
}

func TestAccIdpsDatasource_FilterByName(t *testing.T) {
	frame := test_utils.NewInstanceTestFrame(t, "zitadel_idps")
	googleDep, githubDep := idpDeps(frame)

	config := `
data "zitadel_idps" "default" {
  name       = zitadel_idp_google.default.name
  depends_on = [zitadel_idp_github.default]
}`

	test_utils.RunDatasourceTest(
		t,
		frame.BaseTestFrame,
		config,
		[]string{googleDep, githubDep},
		nil,
		map[string]string{
			"ids.#": "1",
		},
	)
}

func TestAccIdpsDatasource_FilterByType(t *testing.T) {
	frame := test_utils.NewInstanceTestFrame(t, "zitadel_idps")
	googleDep, githubDep := idpDeps(frame)

	config := fmt.Sprintf(`
data "zitadel_idps" "default" {
  name        = "%s"
  name_method = "TEXT_QUERY_METHOD_CONTAINS"
  type        = "PROVIDER_TYPE_GITHUB"
  depends_on  = [zitadel_idp_google.default, zitadel_idp_github.default]
}`, frame.UniqueResourcesID)

	test_utils.RunDatasourceTest(
		t,
		frame.BaseTestFrame,
		config,
		[]string{googleDep, githubDep},
		nil,
		map[string]string{
			"ids.#": "1",
		},
	)
}

func TestAccIdpsDatasource_NoMatch(t *testing.T) {
	frame := test_utils.NewInstanceTestFrame(t, "zitadel_idps")
	googleDep, _ := idpDeps(frame)

	config := `
data "zitadel_idps" "default" {
  name       = "nonexistent"
  depends_on = [zitadel_idp_google.default]
}`

	test_utils.RunDatasourceTest(
		t,
		frame.BaseTestFrame,
		config,
		[]string{googleDep},
		nil,
		map[string]string{
			"ids.#": "0",
		},
	)
}

func idpDeps(frame *test_utils.InstanceTestFrame) (string, string) {
	googleDep := fmt.Sprintf(`
resource "zitadel_idp_google" "default" {
  name                = "google_%s"
  client_id           = "dummy"
  client_secret       = "dummy"
  scopes              = ["openid", "profile", "email"]
  is_linking_allowed  = false
  is_creation_allowed = true
  is_auto_creation    = false
  is_auto_update      = true
}`, frame.UniqueResourcesID)

	githubDep := fmt.Sprintf(`
resource "zitadel_idp_github" "default" {
  name                = "github_%s"
  client_id           = "dummy"
  client_secret       = "dummy"
  scopes              = ["openid", "profile", "email"]
  is_linking_allowed  = false
  is_creation_allowed = true
  is_auto_creation    = false
  is_auto_update      = true
}`, frame.UniqueResourcesID)

	return googleDep, githubDep
}
