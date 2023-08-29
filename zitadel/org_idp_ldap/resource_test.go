package org_idp_ldap_test

import (
	"testing"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/idp_ldap"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/org_idp_utils/org_idp_test_utils"
)

func TestAccOrgIdPLDAP(t *testing.T) {
	org_idp_test_utils.RunOrgLifecyleTest(t, "zitadel_org_idp_ldap", idp_ldap.BindPasswordVar)
}
