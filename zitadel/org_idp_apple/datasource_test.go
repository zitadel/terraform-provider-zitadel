package org_idp_apple_test

import (
	"fmt"
	"testing"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper/test_utils"
)

func TestAccOrgIdpAppleDatasource(t *testing.T) {
	frame := test_utils.NewOrgTestFrame(t, "zitadel_org_idp_apple")
	resourceDep := fmt.Sprintf(`
resource "zitadel_org_idp_apple" "default" {
  org_id              = data.zitadel_org.default.id
  name                = "%s"
  client_id           = "com.example.app"
  team_id             = "ABCDE12345"
  key_id              = "FGHIJ67890"
  private_key         = "dummyprivatekey"
  scopes              = ["name", "email"]
  is_linking_allowed  = false
  is_creation_allowed = true
  is_auto_creation    = false
  is_auto_update      = true
}`, frame.UniqueResourcesID)

	config := `
data "zitadel_org_idp_apple" "default" {
  org_id = data.zitadel_org.default.id
  id     = zitadel_org_idp_apple.default.id
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
