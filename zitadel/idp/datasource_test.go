package idp_test

import (
	"fmt"
	"testing"

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
