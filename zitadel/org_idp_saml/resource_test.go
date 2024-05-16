package org_idp_saml_test

import (
	"testing"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/org_idp_utils/org_idp_test_utils"
)

func TestAccOrgIdPSAML(t *testing.T) {
	org_idp_test_utils.RunOrgLifecyleTest(t, "zitadel_org_idp_saml", "")
}
