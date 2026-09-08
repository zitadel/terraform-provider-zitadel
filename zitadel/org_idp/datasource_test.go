package org_idp_test

import (
	"fmt"
	"testing"

	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/admin"
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

func TestAccOrgIdpsDatasource_FilterByOwnerTypeSystem(t *testing.T) {
	frame := test_utils.NewOrgTestFrame(t, "zitadel_org_idps")
	// An instance-level IDP is visible to the organization with the system owner type.
	name := "instance_google_" + frame.UniqueResourcesID
	resp, err := frame.Admin.AddGoogleProvider(frame, &admin.AddGoogleProviderRequest{
		Name:         name,
		ClientId:     "dummy",
		ClientSecret: "dummy",
	})
	if err != nil {
		t.Fatal(err)
	}
	addGoogleProvider(t, frame, "org_google_"+frame.UniqueResourcesID)

	config := fmt.Sprintf(`
data "zitadel_org_idps" "default" {
  org_id     = "%s"
  name       = "%s"
  owner_type = "IDP_OWNER_TYPE_SYSTEM"
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
			"idps.0.id":         resp.GetId(),
			"idps.0.owner_type": "IDP_OWNER_TYPE_SYSTEM",
		},
	)
}

func TestAccOrgIdpsDatasource_UnspecifiedAppliesNoFilter(t *testing.T) {
	frame := test_utils.NewOrgTestFrame(t, "zitadel_org_idps")
	prefix := "org_idps_" + frame.UniqueResourcesID
	addGoogleProvider(t, frame, prefix+"_google")
	addGitHubProvider(t, frame, prefix+"_github")

	config := fmt.Sprintf(`
data "zitadel_org_idps" "default" {
  org_id      = "%s"
  name        = "%s"
  name_method = "TEXT_QUERY_METHOD_STARTS_WITH"
  type        = "PROVIDER_TYPE_UNSPECIFIED"
  owner_type  = "IDP_OWNER_TYPE_UNSPECIFIED"
}
`, frame.OrgID, prefix)

	test_utils.RunDatasourceTest(
		t,
		frame.BaseTestFrame,
		config,
		[]string{frame.AsOrgDefaultDependency},
		nil,
		map[string]string{
			"idps.#": "2",
		},
	)
}

func TestAccOrgIdpsDatasource_MoreThanOnePage(t *testing.T) {
	frame := test_utils.NewOrgTestFrame(t, "zitadel_org_idps")
	// One more IDP than a single ListProviders page holds, so the datasource must paginate.
	idpCount := 101
	prefix := "org_idps_page_" + frame.UniqueResourcesID
	for i := 0; i < idpCount; i++ {
		addGitHubProvider(t, frame, fmt.Sprintf("%s_%03d", prefix, i))
	}

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
			"idps.#": fmt.Sprint(idpCount),
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
