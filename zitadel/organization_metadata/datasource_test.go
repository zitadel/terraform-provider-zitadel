package organization_metadata_test

import (
	"fmt"
	"testing"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper/test_utils"
)

func TestAccOrganizationMetadataDatasource(t *testing.T) {
	frame := test_utils.NewInstanceTestFrame(t, "zitadel_organization_metadata")
	orgDep := fmt.Sprintf(`
resource "zitadel_organization" "default" {
  name = "%s"
}`, frame.UniqueResourcesID)

	metadataDep := fmt.Sprintf(`
resource "zitadel_organization_metadata" "default" {
  organization_id = zitadel_organization.default.id
  key             = "test_key_%s"
  value           = "test_value"
}`, frame.UniqueResourcesID)

	config := fmt.Sprintf(`
data "zitadel_organization_metadata" "default" {
  organization_id = zitadel_organization.default.id
  key             = "test_key_%s"
  depends_on      = [zitadel_organization_metadata.default]
}`, frame.UniqueResourcesID)

	test_utils.RunDatasourceTest(
		t,
		frame.BaseTestFrame,
		config,
		[]string{orgDep, metadataDep},
		nil,
		map[string]string{
			"value": "test_value",
		},
	)
}

func TestAccOrganizationMetadatasDatasource(t *testing.T) {
	frame := test_utils.NewInstanceTestFrame(t, "zitadel_organization_metadatas")
	orgDep := fmt.Sprintf(`
resource "zitadel_organization" "default" {
  name = "%s"
}`, frame.UniqueResourcesID)

	metadataDep := fmt.Sprintf(`
resource "zitadel_organization_metadata" "default" {
  organization_id = zitadel_organization.default.id
  key             = "test_key_%s"
  value           = "test_value"
}`, frame.UniqueResourcesID)

	config := `
data "zitadel_organization_metadatas" "default" {
  organization_id = zitadel_organization.default.id
  depends_on      = [zitadel_organization_metadata.default]
}`

	test_utils.RunDatasourceTest(
		t,
		frame.BaseTestFrame,
		config,
		[]string{orgDep, metadataDep},
		nil,
		map[string]string{
			"metadata.0.value": "test_value",
		},
	)
}
