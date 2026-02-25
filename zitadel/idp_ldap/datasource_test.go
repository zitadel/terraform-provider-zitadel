package idp_ldap_test

import (
	"fmt"
	"testing"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper/test_utils"
)

func TestAccIdpLdapDatasource(t *testing.T) {
	frame := test_utils.NewInstanceTestFrame(t, "zitadel_idp_ldap")
	resourceDep := fmt.Sprintf(`
resource "zitadel_idp_ldap" "default" {
  name                 = "%s"
  servers              = ["ldap://localhost:389"]
  start_tls            = false
  base_dn              = "dc=example,dc=com"
  bind_dn              = "cn=admin,dc=example,dc=com"
  bind_password        = "password"
  user_base            = "ou=users,dc=example,dc=com"
  user_object_classes  = ["inetOrgPerson"]
  user_filters         = ["uid"]
  timeout              = "10s"
  is_linking_allowed   = false
  is_creation_allowed  = true
  is_auto_creation     = false
  is_auto_update       = true
}`, frame.UniqueResourcesID)

	config := `
data "zitadel_idp_ldap" "default" {
  id = zitadel_idp_ldap.default.id
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
