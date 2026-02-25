package org_idp_jwt_test

import (
	"fmt"
	"testing"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper/test_utils"
)

func TestAccOrgIdpJwtDatasource(t *testing.T) {
	frame := test_utils.NewOrgTestFrame(t, "zitadel_org_idp_jwt")
	resourceDep := fmt.Sprintf(`
resource "zitadel_org_idp_jwt" "default" {
  org_id         = data.zitadel_org.default.id
  name           = "%s"
  styling_type   = "STYLING_TYPE_UNSPECIFIED"
  jwt_endpoint   = "https://jwt.example.com/token"
  keys_endpoint  = "https://jwt.example.com/.well-known/jwks.json"
  issuer         = "https://jwt.example.com"
  header_name    = "x-auth-token"
  auto_register  = false
}`, frame.UniqueResourcesID)

	config := `
data "zitadel_org_idp_jwt" "default" {
  org_id = data.zitadel_org.default.id
  idp_id = zitadel_org_idp_jwt.default.id
}`

	test_utils.RunDatasourceTest(
		t,
		frame.BaseTestFrame,
		config,
		[]string{frame.AsOrgDefaultDependency, resourceDep},
		nil,
		map[string]string{
			"name": frame.UniqueResourcesID,
		},
	)
}
