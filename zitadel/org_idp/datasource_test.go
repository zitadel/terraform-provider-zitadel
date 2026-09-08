package org_idp_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/management"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper/test_utils"
)

func TestAccOrgIdpsDatasource_All(t *testing.T) {
	datasourceName := "zitadel_org_idps"
	frame := test_utils.NewOrgTestFrame(t, datasourceName)

	names := []string{"google_" + frame.UniqueResourcesID, "github_" + frame.UniqueResourcesID}
	_, err := frame.AddGoogleProvider(frame, &management.AddGoogleProviderRequest{
		Name:         names[0],
		ClientId:     "dummy",
		ClientSecret: "dummy",
	})
	if err != nil {
		t.Fatal(err)
	}
	_, err = frame.AddGitHubProvider(frame, &management.AddGitHubProviderRequest{
		Name:         names[1],
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
  name_method = "TEXT_QUERY_METHOD_CONTAINS"
}
`, frame.OrgID, frame.UniqueResourcesID)

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
	_, err := frame.AddGoogleProvider(frame, &management.AddGoogleProviderRequest{
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
		checkIdpExists(frame, matchingName),
		map[string]string{
			"ids.#": "1",
		},
	)
}

func TestAccOrgIdpsDatasource_FilterByType(t *testing.T) {
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
  org_id      = "%s"
  name        = "%s"
  name_method = "TEXT_QUERY_METHOD_CONTAINS"
  type        = "PROVIDER_TYPE_GITHUB"
}
`, frame.OrgID, frame.UniqueResourcesID)

	test_utils.RunDatasourceTest(
		t,
		frame.BaseTestFrame,
		config,
		[]string{frame.AsOrgDefaultDependency},
		nil,
		map[string]string{
			"ids.#": "1",
		},
	)
}

func TestAccOrgIdpsDatasource_FilterByOwnerType(t *testing.T) {
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
  org_id      = "%s"
  name        = "%s"
  name_method = "TEXT_QUERY_METHOD_CONTAINS"
  owner_type  = "IDP_OWNER_TYPE_ORG"
}
`, frame.OrgID, frame.UniqueResourcesID)

	test_utils.RunDatasourceTest(
		t,
		frame.BaseTestFrame,
		config,
		[]string{frame.AsOrgDefaultDependency},
		nil,
		map[string]string{
			"ids.#": "1",
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

func checkIdpExists(frame *test_utils.OrgTestFrame, expectedName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		resp, err := frame.ListProviders(frame, &management.ListProvidersRequest{})
		if err != nil {
			return err
		}

		for _, idp := range resp.Result {
			if idp.Name == expectedName {
				return nil
			}
		}

		return fmt.Errorf("expected idp %s not found", expectedName)
	}
}
