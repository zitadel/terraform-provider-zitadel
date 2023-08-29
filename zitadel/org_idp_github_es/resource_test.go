package org_idp_github_es_test

import (
	"testing"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/idp_utils"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/org_idp_utils/org_idp_test_utils"
)

func TestAccOrgIdPGitHubES(t *testing.T) {
	org_idp_test_utils.RunOrgLifecyleTest(t, "zitadel_org_idp_github_es", idp_utils.ClientSecretVar)
}
