package idp_apple_test

import (
	"fmt"
	"testing"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper/test_utils"
)

func TestAccIdpAppleDatasource(t *testing.T) {
	frame := test_utils.NewInstanceTestFrame(t, "zitadel_idp_apple")
	resourceDep := fmt.Sprintf(`
resource "zitadel_idp_apple" "default" {
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
data "zitadel_idp_apple" "default" {
  id = zitadel_idp_apple.default.id
}`

	test_utils.RunDatasourceTest(
		t,
		frame.BaseTestFrame,
		config,
		[]string{resourceDep},
		nil,
		map[string]string{
			"name": frame.UniqueResourcesID,
		},
	)
}
