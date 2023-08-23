package org_idp_oidc_test

import (
	"testing"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/org_idp_utils/org_idp_test_utils"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/idp_utils"
)

func TestAccOrgIDPOIDC(t *testing.T) {
	org_idp_test_utils.RunOrgLifecyleTest(t, "zitadel_org_idp_oidc", idp_utils.ClientSecretVar)
}
