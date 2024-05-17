package idp_oauth_test

import (
	"testing"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/idp_utils"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/idp_utils/idp_test_utils"
)

func TestAccInstanceIdPOAuth(t *testing.T) {
	idp_test_utils.RunInstanceIDPLifecyleTest(t, "zitadel_idp_oauth", idp_utils.ClientSecretVar)
}
