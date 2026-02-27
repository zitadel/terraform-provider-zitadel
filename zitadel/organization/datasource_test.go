package organization_test

import (
	"fmt"
	"testing"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper/test_utils"
)

func TestAccOrganizationDatasource(t *testing.T) {
	frame := test_utils.NewOrgTestFrame(t, "zitadel_organization")
	resourceDep := fmt.Sprintf(`
resource "zitadel_organization" "default" {
  name = "%s"
}`, frame.UniqueResourcesID)

	config := `
data "zitadel_organization" "default" {
  id = zitadel_organization.default.id
}`

	test_utils.RunDatasourceTest(
		t,
		frame.BaseTestFrame,
		config,
		[]string{resourceDep},
		nil,
		map[string]string{
			"name":  frame.UniqueResourcesID,
			"state": "ORGANIZATION_STATE_ACTIVE",
		},
	)
}

func TestAccOrganizationsDatasource(t *testing.T) {
	frame := test_utils.NewOrgTestFrame(t, "zitadel_organizations")
	config := `
data "zitadel_organizations" "default" {
  state = "ORGANIZATION_STATE_ACTIVE"
}`

	test_utils.RunDatasourceTest(
		t,
		frame.BaseTestFrame,
		config,
		nil,
		nil,
		map[string]string{},
	)
}
