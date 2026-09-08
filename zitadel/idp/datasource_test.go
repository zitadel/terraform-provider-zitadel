package idp_test

import (
	"fmt"
	"testing"

	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/admin"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper/test_utils"
)

func TestAccIdpsDatasource_FilterByName(t *testing.T) {
	datasourceName := "zitadel_idps"
	frame := test_utils.NewInstanceTestFrame(t, datasourceName)

	matchingName := "google_" + frame.UniqueResourcesID
	google, err := frame.AddGoogleProvider(frame, &admin.AddGoogleProviderRequest{
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
		nil,
		map[string]string{
			"ids.#": "1",
			"ids.0": google.GetId(),
		},
	)
}

func TestAccIdpsDatasource_FilterByType(t *testing.T) {
	datasourceName := "zitadel_idps"
	frame := test_utils.NewInstanceTestFrame(t, datasourceName)

	// Both providers share the name prefix, so only the type filter narrows the result down.
	prefix := "idps_" + frame.UniqueResourcesID
	_, err := frame.AddGoogleProvider(frame, &admin.AddGoogleProviderRequest{
		Name:         prefix + "_google",
		ClientId:     "dummy",
		ClientSecret: "dummy",
	})
	if err != nil {
		t.Fatal(err)
	}
	github, err := frame.AddGitHubProvider(frame, &admin.AddGitHubProviderRequest{
		Name:         prefix + "_github",
		ClientId:     "dummy",
		ClientSecret: "dummy",
	})
	if err != nil {
		t.Fatal(err)
	}

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
			"ids.0": github.GetId(),
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
