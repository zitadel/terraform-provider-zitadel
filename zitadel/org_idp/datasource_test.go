package org_idp_test

import (
	"fmt"
	"testing"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper/test_utils"
)

func TestAccOrgIdpsDatasource_All(t *testing.T) {
	frame := test_utils.NewOrgTestFrame(t, "zitadel_org_idps")
	googleDep, githubDep := orgIdpDeps(frame)

	config := fmt.Sprintf(`
data "zitadel_org_idps" "default" {
  org_id      = data.zitadel_org.default.id
  name        = "%s"
  name_method = "TEXT_QUERY_METHOD_CONTAINS"
  depends_on  = [zitadel_org_idp_google.default, zitadel_org_idp_github.default]
}`, frame.UniqueResourcesID)

	test_utils.RunDatasourceTest(
		t,
		frame.BaseTestFrame,
		config,
		[]string{frame.AsOrgDefaultDependency, googleDep, githubDep},
		nil,
		map[string]string{
			"ids.#": "2",
		},
	)
}

func TestAccOrgIdpsDatasource_FilterByName(t *testing.T) {
	frame := test_utils.NewOrgTestFrame(t, "zitadel_org_idps")
	googleDep, githubDep := orgIdpDeps(frame)

	config := `
data "zitadel_org_idps" "default" {
  org_id     = data.zitadel_org.default.id
  name       = zitadel_org_idp_google.default.name
  depends_on = [zitadel_org_idp_github.default]
}`

	test_utils.RunDatasourceTest(
		t,
		frame.BaseTestFrame,
		config,
		[]string{frame.AsOrgDefaultDependency, googleDep, githubDep},
		nil,
		map[string]string{
			"ids.#": "1",
		},
	)
}

func TestAccOrgIdpsDatasource_FilterByType(t *testing.T) {
	frame := test_utils.NewOrgTestFrame(t, "zitadel_org_idps")
	googleDep, githubDep := orgIdpDeps(frame)

	config := fmt.Sprintf(`
data "zitadel_org_idps" "default" {
  org_id      = data.zitadel_org.default.id
  name        = "%s"
  name_method = "TEXT_QUERY_METHOD_CONTAINS"
  type        = "PROVIDER_TYPE_GITHUB"
  depends_on  = [zitadel_org_idp_google.default, zitadel_org_idp_github.default]
}`, frame.UniqueResourcesID)

	test_utils.RunDatasourceTest(
		t,
		frame.BaseTestFrame,
		config,
		[]string{frame.AsOrgDefaultDependency, googleDep, githubDep},
		nil,
		map[string]string{
			"ids.#": "1",
		},
	)
}

func TestAccOrgIdpsDatasource_FilterByOwnerType(t *testing.T) {
	frame := test_utils.NewOrgTestFrame(t, "zitadel_org_idps")
	googleDep, _ := orgIdpDeps(frame)

	instanceDep := fmt.Sprintf(`
resource "zitadel_idp_google" "default" {
  name                = "instance_google_%s"
  client_id           = "dummy"
  client_secret       = "dummy"
  scopes              = ["openid", "profile", "email"]
  is_linking_allowed  = false
  is_creation_allowed = true
  is_auto_creation    = false
  is_auto_update      = true
}`, frame.UniqueResourcesID)

	config := fmt.Sprintf(`
data "zitadel_org_idps" "default" {
  org_id      = data.zitadel_org.default.id
  name        = "%s"
  name_method = "TEXT_QUERY_METHOD_CONTAINS"
  owner_type  = "IDP_OWNER_TYPE_ORG"
  depends_on  = [zitadel_org_idp_google.default, zitadel_idp_google.default]
}`, frame.UniqueResourcesID)

	test_utils.RunDatasourceTest(
		t,
		frame.BaseTestFrame,
		config,
		[]string{frame.AsOrgDefaultDependency, googleDep, instanceDep},
		nil,
		map[string]string{
			"ids.#": "1",
		},
	)
}

func TestAccOrgIdpsDatasource_NoMatch(t *testing.T) {
	frame := test_utils.NewOrgTestFrame(t, "zitadel_org_idps")
	googleDep, _ := orgIdpDeps(frame)

	config := `
data "zitadel_org_idps" "default" {
  org_id     = data.zitadel_org.default.id
  name       = "nonexistent"
  depends_on = [zitadel_org_idp_google.default]
}`

	test_utils.RunDatasourceTest(
		t,
		frame.BaseTestFrame,
		config,
		[]string{frame.AsOrgDefaultDependency, googleDep},
		nil,
		map[string]string{
			"ids.#": "0",
		},
	)
}

func orgIdpDeps(frame *test_utils.OrgTestFrame) (string, string) {
	googleDep := fmt.Sprintf(`
resource "zitadel_org_idp_google" "default" {
  org_id              = data.zitadel_org.default.id
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
resource "zitadel_org_idp_github" "default" {
  org_id              = data.zitadel_org.default.id
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
