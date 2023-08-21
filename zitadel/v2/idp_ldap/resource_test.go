package idp_ldap_test

import (
	"fmt"
	"testing"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/helper/test_utils"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/idp_ldap"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/idp_utils/idp_test_utils"
)

func TestAccInstanceIdPLDAP(t *testing.T) {
	resourceName := "zitadel_idp_ldap"
	frame, err := test_utils.NewInstanceTestFrame(resourceName)
	if err != nil {
		t.Fatalf("setting up test context failed: %v", err)
	}
	idp_test_utils.RunInstanceIDPLifecyleTest(t, *frame, func(name, secret string) string {
		return fmt.Sprintf(`
resource "%s" "%s" {
  name                  = "%s"
  servers               = ["a server"]
  start_tls             = true
  base_dn               = "a base dn"
  bind_dn               = "a bind dn"
  bind_password         = "%s"
  user_base             = "a user base"
  user_object_classes   = ["a user object class"]
  user_filters          = ["a user filter"]
  timeout               = "5s"
  id_attribute          = "a id_attribute"
  first_name_attribute  = "a first name attribute"
  last_name_attribute   = "a last name attribute"
  is_linking_allowed    = false
  is_creation_allowed   = true
  is_auto_creation      = false
  is_auto_update        = true
}`, resourceName, frame.UniqueResourcesID, name, secret)
	}, idp_ldap.BindPasswordVar)
}
