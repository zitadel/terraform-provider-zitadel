package idp_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/admin"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper/test_utils"
)

func TestAccIdpsDatasource_All(t *testing.T) {
	datasourceName := "zitadel_idps"
	frame := test_utils.NewInstanceTestFrame(t, datasourceName)

	names := []string{"google_" + frame.UniqueResourcesID, "github_" + frame.UniqueResourcesID}
	_, err := frame.AddGoogleProvider(frame, &admin.AddGoogleProviderRequest{
		Name:         names[0],
		ClientId:     "dummy",
		ClientSecret: "dummy",
	})
	if err != nil {
		t.Fatal(err)
	}
	_, err = frame.AddGitHubProvider(frame, &admin.AddGitHubProviderRequest{
		Name:         names[1],
		ClientId:     "dummy",
		ClientSecret: "dummy",
	})
	if err != nil {
		t.Fatal(err)
	}

	config := fmt.Sprintf(`
data "zitadel_idps" "default" {
  name        = "%s"
  name_method = "TEXT_QUERY_METHOD_CONTAINS"
}
`, frame.UniqueResourcesID)

	test_utils.RunDatasourceTest(
		t,
		frame.BaseTestFrame,
		config,
		nil,
		nil,
		map[string]string{
			"ids.#": "2",
		},
	)
}

func TestAccIdpsDatasource_FilterByName(t *testing.T) {
	datasourceName := "zitadel_idps"
	frame := test_utils.NewInstanceTestFrame(t, datasourceName)

	matchingName := "google_" + frame.UniqueResourcesID
	_, err := frame.AddGoogleProvider(frame, &admin.AddGoogleProviderRequest{
		Name:         matchingName,
		ClientId:     "dummy",
		ClientSecret: "dummy",
	})
	if err != nil {
		t.Fatal(err)
	}
	_, err = frame.AddGitHubProvider(frame, &admin.AddGitHubProviderRequest{
		Name:         "github_" + frame.UniqueResourcesID,
		ClientId:     "dummy",
		ClientSecret: "dummy",
	})
	if err != nil {
		t.Fatal(err)
	}

	config := fmt.Sprintf(`
data "zitadel_idps" "default" {
  name = "%s"
}
`, matchingName)

	test_utils.RunDatasourceTest(
		t,
		frame.BaseTestFrame,
		config,
		nil,
		checkIdpExists(frame, matchingName),
		map[string]string{
			"ids.#": "1",
		},
	)
}

func TestAccIdpsDatasource_FilterByType(t *testing.T) {
	datasourceName := "zitadel_idps"
	frame := test_utils.NewInstanceTestFrame(t, datasourceName)

	_, err := frame.AddGoogleProvider(frame, &admin.AddGoogleProviderRequest{
		Name:         "google_" + frame.UniqueResourcesID,
		ClientId:     "dummy",
		ClientSecret: "dummy",
	})
	if err != nil {
		t.Fatal(err)
	}
	_, err = frame.AddGitHubProvider(frame, &admin.AddGitHubProviderRequest{
		Name:         "github_" + frame.UniqueResourcesID,
		ClientId:     "dummy",
		ClientSecret: "dummy",
	})
	if err != nil {
		t.Fatal(err)
	}

	config := fmt.Sprintf(`
data "zitadel_idps" "default" {
  name        = "%s"
  name_method = "TEXT_QUERY_METHOD_CONTAINS"
  type        = "PROVIDER_TYPE_GITHUB"
}
`, frame.UniqueResourcesID)

	test_utils.RunDatasourceTest(
		t,
		frame.BaseTestFrame,
		config,
		nil,
		nil,
		map[string]string{
			"ids.#": "1",
		},
	)
}

func TestAccIdpsDatasource_NoMatch(t *testing.T) {
	datasourceName := "zitadel_idps"
	frame := test_utils.NewInstanceTestFrame(t, datasourceName)

	_, err := frame.AddGoogleProvider(frame, &admin.AddGoogleProviderRequest{
		Name:         "google_" + frame.UniqueResourcesID,
		ClientId:     "dummy",
		ClientSecret: "dummy",
	})
	if err != nil {
		t.Fatal(err)
	}

	config := `
data "zitadel_idps" "default" {
  name = "nonexistent"
}
`

	test_utils.RunDatasourceTest(
		t,
		frame.BaseTestFrame,
		config,
		nil,
		nil,
		map[string]string{
			"ids.#": "0",
		},
	)
}

func checkIdpExists(frame *test_utils.InstanceTestFrame, expectedName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		resp, err := frame.ListProviders(frame, &admin.ListProvidersRequest{})
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
