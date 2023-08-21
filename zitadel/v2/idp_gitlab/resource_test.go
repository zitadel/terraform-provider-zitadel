package idp_gitlab_test

import (
	"testing"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/idp_utils"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/idp_utils/idp_test_utils"
)

func TestAccInstanceIdPGitLab(t *testing.T) {
	idp_test_utils.RunInstanceIDPLifecyleTest(t, "zitadel_idp_gitlab", idp_utils.ClientSecretVar)
}
