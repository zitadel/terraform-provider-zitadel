package org_idp_azure_ad_test

import (
	"testing"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/idp_utils"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/org_idp_utils/org_idp_test_utils"
)

func TestAccOrgIdPAzureAD(t *testing.T) {
	org_idp_test_utils.RunOrgLifecyleTest(t, "zitadel_org_idp_azure_ad", idp_utils.ClientSecretVar)
}
