package idp_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/admin"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper/test_utils"
)

func TestAccIdpsDatasource_FilterByName(t *testing.T) {
	frame := test_utils.NewInstanceTestFrame(t, "zitadel_idps")
	name := "google_" + frame.UniqueResourcesID
	googleID := addGoogleProvider(t, frame, name)
	addGitHubProvider(t, frame, "github_"+frame.UniqueResourcesID)

	config := fmt.Sprintf(`
data "zitadel_idps" "default" {
  name = "%s"
}
`, name)

	test_utils.RunDatasourceTest(
		t,
		frame.BaseTestFrame,
		config,
		nil,
		nil,
		map[string]string{
			"idps.#":            "1",
			"idps.0.id":         googleID,
			"idps.0.name":       name,
			"idps.0.type":       "PROVIDER_TYPE_GOOGLE",
			"idps.0.state":      "IDP_STATE_ACTIVE",
			"idps.0.owner_type": "IDP_OWNER_TYPE_SYSTEM",
		},
	)
}

func TestAccIdpsDatasource_FilterByNameAndType(t *testing.T) {
	frame := test_utils.NewInstanceTestFrame(t, "zitadel_idps")
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
			"idps.#":      "1",
			"idps.0.id":   githubID,
			"idps.0.type": "PROVIDER_TYPE_GITHUB",
		},
	)
}

func TestAccIdpsDatasource_FilterByTypeOnly(t *testing.T) {
	frame := test_utils.NewInstanceTestFrame(t, "zitadel_idps")
	addGoogleProvider(t, frame, "google_"+frame.UniqueResourcesID)
	githubID := addGitHubProvider(t, frame, "github_"+frame.UniqueResourcesID)

	config := `
data "zitadel_idps" "default" {
  type = "PROVIDER_TYPE_GITHUB"
}
`

	// Other tests create instance IDPs concurrently, so only assert the type of every
	// element and that the GitHub IDP created here is included.
	test_utils.RunDatasourceTest(
		t,
		frame.BaseTestFrame,
		config,
		nil,
		checkAllOfTypeAndContains(frame, "PROVIDER_TYPE_GITHUB", githubID),
		nil,
	)
}

func TestAccIdpsDatasource_NoMatch(t *testing.T) {
	frame := test_utils.NewInstanceTestFrame(t, "zitadel_idps")
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
			"idps.#": "0",
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

func checkAllOfTypeAndContains(frame *test_utils.InstanceTestFrame, expectedType, expectedID string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		attrs := frame.State(state).Attributes
		count, err := strconv.Atoi(attrs["idps.#"])
		if err != nil {
			return fmt.Errorf("idps.# is not a number: %v", err)
		}
		found := false
		for i := 0; i < count; i++ {
			if actualType := attrs[fmt.Sprintf("idps.%d.type", i)]; actualType != expectedType {
				return fmt.Errorf("expected only idps of type %s, but idps.%d has type %s", expectedType, i, actualType)
			}
			if attrs[fmt.Sprintf("idps.%d.id", i)] == expectedID {
				found = true
			}
		}
		if !found {
			return fmt.Errorf("expected idp %s to be listed", expectedID)
		}
		return nil
	}
}
