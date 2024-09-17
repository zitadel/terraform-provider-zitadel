package idp_google_test

import (
	"testing"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/idp_utils"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/idp_utils/idp_test_utils"
)

func TestAccInstanceIdPGoogle(t *testing.T) {
	idp_test_utils.RunInstanceIDPLifecyleTest(t, "zitadel_idp_google", idp_utils.ClientSecretVar)
}
