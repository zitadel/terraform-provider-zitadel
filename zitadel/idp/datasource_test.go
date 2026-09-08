package idp_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/admin"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper/test_utils"
)

func TestAccIdpsDatasource_FilterByName(t *testing.T) {
	datasourceName := "zitadel_idps"
	frame := test_utils.NewInstanceTestFrame(t, datasourceName)

	matchingName := "google_" + frame.UniqueResourcesID
	googleID := addGoogleProvider(t, frame, matchingName)
	addGitHubProvider(t, frame, "github_"+frame.UniqueResourcesID)

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
		nil,
		map[string]string{
			"ids.#": "1",
			"ids.0": googleID,
		},
	)
}

func TestAccIdpsDatasource_FilterByNameAndType(t *testing.T) {
	datasourceName := "zitadel_idps"
	frame := test_utils.NewInstanceTestFrame(t, datasourceName)

	// Both providers share the name prefix, so only the type filter narrows the result down.
	prefix := "idps_" + frame.UniqueResourcesID
	addGoogleProvider(t, frame, prefix+"_google")
	githubID := addGitHubProvider(t, frame, prefix+"_github")

	config := fmt.Sprintf(`
data "zitadel_idps" "default" {
  name        = "%s"
  name_method = "TEXT_QUERY_METHOD_STARTS_WITH"
  type        = "PROVIDER_TYPE_GITHUB"
}
`, prefix)

	test_utils.RunDatasourceTest(
		t,
		frame.BaseTestFrame,
		config,
		nil,
		nil,
		map[string]string{
			"ids.#": "1",
			"ids.0": githubID,
		},
	)
}

func TestAccIdpsDatasource_FilterByType(t *testing.T) {
	datasourceName := "zitadel_idps"
	frame := test_utils.NewInstanceTestFrame(t, datasourceName)

	googleID := addGoogleProvider(t, frame, "google_"+frame.UniqueResourcesID)
	githubID := addGitHubProvider(t, frame, "github_"+frame.UniqueResourcesID)

	config := `
data "zitadel_idps" "default" {
  type = "PROVIDER_TYPE_GITHUB"
}
`

	// Other tests create instance IDPs concurrently, so only check the presence of the IDs created here.
	test_utils.RunDatasourceTest(
		t,
		frame.BaseTestFrame,
		config,
		nil,
		checkIDsContain(frame, githubID, googleID),
		nil,
	)
}

func TestAccIdpsDatasource_NoMatch(t *testing.T) {
	datasourceName := "zitadel_idps"
	frame := test_utils.NewInstanceTestFrame(t, datasourceName)

	addGoogleProvider(t, frame, "google_"+frame.UniqueResourcesID)

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

func addGoogleProvider(t *testing.T, frame *test_utils.InstanceTestFrame, name string) string {
	resp, err := frame.AddGoogleProvider(frame, &admin.AddGoogleProviderRequest{
		Name:         name,
		ClientId:     "dummy",
		ClientSecret: "dummy",
	})
	if err != nil {
		t.Fatal(err)
	}
	return resp.GetId()
}

func addGitHubProvider(t *testing.T, frame *test_utils.InstanceTestFrame, name string) string {
	resp, err := frame.AddGitHubProvider(frame, &admin.AddGitHubProviderRequest{
		Name:         name,
		ClientId:     "dummy",
		ClientSecret: "dummy",
	})
	if err != nil {
		t.Fatal(err)
	}
	return resp.GetId()
}

func checkIDsContain(frame *test_utils.InstanceTestFrame, expectedID, unexpectedID string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		attrs := frame.State(state).Attributes
		found := false
		for key, value := range attrs {
			if key == "ids.#" || len(key) < 4 || key[:4] != "ids." {
				continue
			}
			if value == unexpectedID {
				return fmt.Errorf("expected idp %s not to be listed", unexpectedID)
			}
			if value == expectedID {
				found = true
			}
		}
		if !found {
			return fmt.Errorf("expected idp %s to be listed", expectedID)
		}
		return nil
	}
}
