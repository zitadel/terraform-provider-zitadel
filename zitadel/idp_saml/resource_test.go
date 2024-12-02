package idp_saml_test

import (
	"testing"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/idp_utils/idp_test_utils"
)

func TestAccInstanceIdPSAML(t *testing.T) {
	idp_test_utils.RunInstanceIDPLifecyleTest(t, "zitadel_idp_saml", "")
}
