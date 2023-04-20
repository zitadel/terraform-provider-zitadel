package org_idp_ldap_test

import (
	"fmt"
	"testing"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/helper/test_utils"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/idp_utils"

	test_utils_org "github.com/zitadel/terraform-provider-zitadel/zitadel/v2/org_idp_utils/test_utils"
)

func TestAccZITADELOrgIdPLDAP(t *testing.T) {
	resourceName := "zitadel_org_idp_ldap"
	frame, err := test_utils.NewOrgTestFrame(resourceName)
	if err != nil {
		t.Fatalf("setting up test context failed: %v", err)
	}
	test_utils_org.RunOrgLifecyleTest(t, *frame, func(name, secret string) string {
		return fmt.Sprintf(`
resource "%s" "%s" {
  org_id              = "%s"
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
}`, resourceName, frame.UniqueResourcesID, frame.OrgID, name, secret)
	}, idp_utils.BindPasswordVar)
}
