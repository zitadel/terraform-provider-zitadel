package idp_azure_ad_test

import (
	"fmt"
	"testing"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper/test_utils"
)

func TestAccIdpAzureADDatasource(t *testing.T) {
	frame := test_utils.NewInstanceTestFrame(t, "zitadel_idp_azure_ad")
	resourceDep := fmt.Sprintf(`
resource "zitadel_idp_azure_ad" "default" {
  name                = "%s"
  client_id           = "dummy"
  client_secret       = "dummy"
  scopes              = ["openid", "profile", "email"]
  tenant_type         = "AZURE_AD_TENANT_TYPE_COMMON"
  email_verified      = true
  is_linking_allowed  = false
  is_creation_allowed = true
  is_auto_creation    = false
  is_auto_update      = true
}`, frame.UniqueResourcesID)

	config := `
data "zitadel_idp_azure_ad" "default" {
  id = zitadel_idp_azure_ad.default.id
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
