package org_idp_test

import (
	"fmt"
	"testing"

	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/management"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper/test_utils"
)

func TestAccOrgIdpsDatasource_FilterByName(t *testing.T) {
	frame := test_utils.NewOrgTestFrame(t, "zitadel_org_idps")
	name := "google_" + frame.UniqueResourcesID
	googleID := addGoogleProvider(t, frame, name)
	addGitHubProvider(t, frame, "github_"+frame.UniqueResourcesID)

	config := fmt.Sprintf(`
data "zitadel_org_idps" "default" {
  org_id = "%s"
  name   = "%s"
}
`, frame.OrgID, name)

	test_utils.RunDatasourceTest(
		t,
		frame.BaseTestFrame,
		config,
		[]string{frame.AsOrgDefaultDependency},
		nil,
		map[string]string{
			"idps.#":            "1",
			"idps.0.id":         googleID,
			"idps.0.name":       name,
			"idps.0.type":       "PROVIDER_TYPE_GOOGLE",
			"idps.0.state":      "IDP_STATE_ACTIVE",
			"idps.0.owner_type": "IDP_OWNER_TYPE_ORG",
		},
	)
}

func TestAccOrgIdpsDatasource_FilterByTypeAndOwner(t *testing.T) {
	frame := test_utils.NewOrgTestFrame(t, "zitadel_org_idps")
	prefix := "org_idps_" + frame.UniqueResourcesID
	addGoogleProvider(t, frame, prefix+"_google")
	githubID := addGitHubProvider(t, frame, prefix+"_github")

	config := fmt.Sprintf(`
data "zitadel_org_idps" "default" {
  org_id      = "%s"
  name        = "%s"
  name_method = "TEXT_QUERY_METHOD_STARTS_WITH"
  type        = "PROVIDER_TYPE_GITHUB"
  owner_type  = "IDP_OWNER_TYPE_ORG"
}
`, frame.OrgID, prefix)

	test_utils.RunDatasourceTest(
		t,
		frame.BaseTestFrame,
		config,
		[]string{frame.AsOrgDefaultDependency},
		nil,
		map[string]string{
			"idps.#":            "1",
			"idps.0.id":         githubID,
			"idps.0.type":       "PROVIDER_TYPE_GITHUB",
			"idps.0.owner_type": "IDP_OWNER_TYPE_ORG",
		},
	)
}

func TestAccOrgIdpsDatasource_NoMatch(t *testing.T) {
	frame := test_utils.NewOrgTestFrame(t, "zitadel_org_idps")
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
			"idps.#": "0",
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
