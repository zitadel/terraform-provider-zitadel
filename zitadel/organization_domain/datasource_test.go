package organization_domain_test

import (
	"fmt"
	"testing"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper/test_utils"
)

func TestAccOrganizationDomainDatasource(t *testing.T) {
	frame := test_utils.NewInstanceTestFrame(t, "zitadel_organization_domain")
	orgDep := fmt.Sprintf(`
resource "zitadel_organization" "default" {
  name = "%s"
}`, frame.UniqueResourcesID)

	domainDep := fmt.Sprintf(`
resource "zitadel_organization_domain" "default" {
  organization_id = zitadel_organization.default.id
  domain          = "%s.example.com"
  validation_type = "DOMAIN_VALIDATION_TYPE_HTTP"
}`, frame.UniqueResourcesID)

	config := fmt.Sprintf(`
data "zitadel_organization_domain" "default" {
  organization_id = zitadel_organization.default.id
  domain          = "%s.example.com"
  depends_on      = [zitadel_organization_domain.default]
}`, frame.UniqueResourcesID)

	test_utils.RunDatasourceTest(
		t,
		frame.BaseTestFrame,
		config,
		[]string{orgDep, domainDep},
		nil,
		map[string]string{},
	)
}

func TestAccOrganizationDomainsDatasource(t *testing.T) {
	frame := test_utils.NewInstanceTestFrame(t, "zitadel_organization_domains")
	orgDep := fmt.Sprintf(`
resource "zitadel_organization" "default" {
  name = "%s"
}`, frame.UniqueResourcesID)

	config := `
data "zitadel_organization_domains" "default" {
  organization_id = zitadel_organization.default.id
}`

	test_utils.RunDatasourceTest(
		t,
		frame.BaseTestFrame,
		config,
		[]string{orgDep},
		nil,
		map[string]string{},
	)
}
