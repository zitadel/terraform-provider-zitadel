package org_idp_gitlab_self_hosted_test

import (
	"testing"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/idp_utils"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/org_idp_utils/org_idp_test_utils"
)

func TestAccOrgIdPGitLabSelfHosted(t *testing.T) {
	org_idp_test_utils.RunOrgLifecyleTest(t, "zitadel_org_idp_gitlab_self_hosted", idp_utils.ClientSecretVar)
}
