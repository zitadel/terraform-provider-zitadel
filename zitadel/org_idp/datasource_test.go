package org_idp_test

import (
	"fmt"
	"testing"

	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/admin"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/management"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper/test_utils"
)

func TestAccOrgIdpsDatasource_All(t *testing.T) {
	datasourceName := "zitadel_org_idps"
	frame := test_utils.NewOrgTestFrame(t, datasourceName)

	prefix := "org_idps_" + frame.UniqueResourcesID
	addGoogleProvider(t, frame, prefix+"_google")
	addGitHubProvider(t, frame, prefix+"_github")

	config := fmt.Sprintf(`
data "zitadel_org_idps" "default" {
  org_id      = "%s"
  name        = "%s"
  name_method = "TEXT_QUERY_METHOD_STARTS_WITH"
}
`, frame.OrgID, prefix)

	test_utils.RunDatasourceTest(
		t,
		frame.BaseTestFrame,
		config,
		[]string{frame.AsOrgDefaultDependency},
		nil,
		map[string]string{
			"ids.#": "2",
		},
	)
}

func TestAccOrgIdpsDatasource_FilterByName(t *testing.T) {
	datasourceName := "zitadel_org_idps"
	frame := test_utils.NewOrgTestFrame(t, datasourceName)

	matchingName := "google_" + frame.UniqueResourcesID
	googleID := addGoogleProvider(t, frame, matchingName)
	addGitHubProvider(t, frame, "github_"+frame.UniqueResourcesID)

	config := fmt.Sprintf(`
data "zitadel_org_idps" "default" {
  org_id = "%s"
  name   = "%s"
}
`, frame.OrgID, matchingName)

	test_utils.RunDatasourceTest(
		t,
		frame.BaseTestFrame,
		config,
		[]string{frame.AsOrgDefaultDependency},
		nil,
		map[string]string{
			"ids.#": "1",
			"ids.0": googleID,
		},
	)
}

func TestAccOrgIdpsDatasource_FilterByType(t *testing.T) {
	datasourceName := "zitadel_org_idps"
	frame := test_utils.NewOrgTestFrame(t, datasourceName)

	prefix := "org_idps_" + frame.UniqueResourcesID
	addGoogleProvider(t, frame, prefix+"_google")
	githubID := addGitHubProvider(t, frame, prefix+"_github")

	config := fmt.Sprintf(`
data "zitadel_org_idps" "default" {
  org_id      = "%s"
  name        = "%s"
  name_method = "TEXT_QUERY_METHOD_STARTS_WITH"
  type        = "PROVIDER_TYPE_GITHUB"
}
`, frame.OrgID, prefix)

	test_utils.RunDatasourceTest(
		t,
		frame.BaseTestFrame,
		config,
		[]string{frame.AsOrgDefaultDependency},
		nil,
		map[string]string{
			"ids.#": "1",
			"ids.0": githubID,
		},
	)
}

func TestAccOrgIdpsDatasource_FilterByOwnerType(t *testing.T) {
	datasourceName := "zitadel_org_idps"
	frame := test_utils.NewOrgTestFrame(t, datasourceName)

	// An instance-level IDP is visible to the organization with the system owner type.
	prefix := "org_idps_" + frame.UniqueResourcesID
	instanceIDP, err := frame.Admin.AddGoogleProvider(frame, &admin.AddGoogleProviderRequest{
		Name:         prefix + "_instance",
		ClientId:     "dummy",
		ClientSecret: "dummy",
	})
	if err != nil {
		t.Fatal(err)
	}
	addGoogleProvider(t, frame, prefix+"_org")

	config := fmt.Sprintf(`
data "zitadel_org_idps" "default" {
  org_id      = "%s"
  name        = "%s"
  name_method = "TEXT_QUERY_METHOD_STARTS_WITH"
  owner_type  = "IDP_OWNER_TYPE_SYSTEM"
}
`, frame.OrgID, prefix)

	test_utils.RunDatasourceTest(
		t,
		frame.BaseTestFrame,
		config,
		[]string{frame.AsOrgDefaultDependency},
		nil,
		map[string]string{
			"ids.#": "1",
			"ids.0": instanceIDP.GetId(),
		},
	)
}

func TestAccOrgIdpsDatasource_NoMatch(t *testing.T) {
	datasourceName := "zitadel_org_idps"
	frame := test_utils.NewOrgTestFrame(t, datasourceName)

	addGoogleProvider(t, frame, "google_"+frame.UniqueResourcesID)

	config := fmt.Sprintf(`
data "zitadel_org_idps" "default" {
  org_id = "%s"
  name   = "nonexistent"
}
`, frame.OrgID)

	test_utils.RunDatasourceTest(
		t,
		frame.BaseTestFrame,
		config,
		[]string{frame.AsOrgDefaultDependency},
		nil,
		map[string]string{
			"ids.#": "0",
		},
	)
}

func addGoogleProvider(t *testing.T, frame *test_utils.OrgTestFrame, name string) string {
	resp, err := frame.AddGoogleProvider(frame, &management.AddGoogleProviderRequest{
		Name:         name,
		ClientId:     "dummy",
		ClientSecret: "dummy",
	})
	if err != nil {
		t.Fatal(err)
	}
	return resp.GetId()
}

func addGitHubProvider(t *testing.T, frame *test_utils.OrgTestFrame, name string) string {
	resp, err := frame.AddGitHubProvider(frame, &management.AddGitHubProviderRequest{
		Name:         name,
		ClientId:     "dummy",
		ClientSecret: "dummy",
	})
	if err != nil {
		t.Fatal(err)
	}
	return resp.GetId()
}
