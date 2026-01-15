package idp_apple_test

import (
	"testing"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/idp_apple"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/idp_utils/idp_test_utils"
)

func TestAccInstanceIdPApple(t *testing.T) {
	idp_test_utils.RunInstanceIDPLifecyleTest(t, "zitadel_idp_apple", idp_apple.PrivateKeyVar)
}
