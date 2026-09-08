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
	_, err := frame.AddGoogleProvider(frame, &management.AddGoogleProviderRequest{
		Name:         prefix + "_google",
		ClientId:     "dummy",
		ClientSecret: "dummy",
	})
	if err != nil {
		t.Fatal(err)
	}
	_, err = frame.AddGitHubProvider(frame, &management.AddGitHubProviderRequest{
		Name:         prefix + "_github",
		ClientId:     "dummy",
		ClientSecret: "dummy",
	})
	if err != nil {
		t.Fatal(err)
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
			"ids.#": "2",
		},
	)
}

func TestAccOrgIdpsDatasource_FilterByName(t *testing.T) {
	datasourceName := "zitadel_org_idps"
	frame := test_utils.NewOrgTestFrame(t, datasourceName)

	matchingName := "google_" + frame.UniqueResourcesID
	google, err := frame.AddGoogleProvider(frame, &management.AddGoogleProviderRequest{
		Name:         matchingName,
		ClientId:     "dummy",
		ClientSecret: "dummy",
	})
	if err != nil {
		t.Fatal(err)
	}
	_, err = frame.AddGitHubProvider(frame, &management.AddGitHubProviderRequest{
		Name:         "github_" + frame.UniqueResourcesID,
		ClientId:     "dummy",
		ClientSecret: "dummy",
	})
	if err != nil {
		t.Fatal(err)
	}

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
			"ids.0": google.GetId(),
		},
	)
}

func TestAccOrgIdpsDatasource_FilterByType(t *testing.T) {
	datasourceName := "zitadel_org_idps"
	frame := test_utils.NewOrgTestFrame(t, datasourceName)

	prefix := "org_idps_" + frame.UniqueResourcesID
	_, err := frame.AddGoogleProvider(frame, &management.AddGoogleProviderRequest{
		Name:         prefix + "_google",
		ClientId:     "dummy",
		ClientSecret: "dummy",
	})
	if err != nil {
		t.Fatal(err)
	}
	github, err := frame.AddGitHubProvider(frame, &management.AddGitHubProviderRequest{
		Name:         prefix + "_github",
		ClientId:     "dummy",
		ClientSecret: "dummy",
	})
	if err != nil {
		t.Fatal(err)
	}

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
			"ids.0": github.GetId(),
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
	_, err = frame.AddGoogleProvider(frame, &management.AddGoogleProviderRequest{
		Name:         prefix + "_org",
		ClientId:     "dummy",
		ClientSecret: "dummy",
	})
	if err != nil {
		t.Fatal(err)
	}

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

	_, err := frame.AddGoogleProvider(frame, &management.AddGoogleProviderRequest{
		Name:         "google_" + frame.UniqueResourcesID,
		ClientId:     "dummy",
		ClientSecret: "dummy",
	})
	if err != nil {
		t.Fatal(err)
	}

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
