package org_idp_oidc_test

import (
	"fmt"
	"testing"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/helper/test_utils"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/idp_utils"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/org_idp_utils/org_idp_test_utils"
)

func TestAccOrgIdPOIDC(t *testing.T) {
	resourceName := "zitadel_org_idp_oidc"
	frame, err := test_utils.NewOrgTestFrame(resourceName)
	if err != nil {
		t.Fatalf("setting up test context failed: %v", err)
	}
	org_idp_test_utils.RunOrgLifecyleTest(t, *frame, func(name, secret string) string {
		return fmt.Sprintf(`
resource "%s" "%s" {
  org_id              = "%s"
  name                = "%s"
  client_id           = "aclientid"
  client_secret       = "%s"
  styling_type         = "STYLING_TYPE_UNSPECIFIED"
  issuer               = "https://google.com"
  scopes               = ["openid", "profile", "email"]
  display_name_mapping = "OIDC_MAPPING_FIELD_PREFERRED_USERNAME"
  username_mapping     = "OIDC_MAPPING_FIELD_PREFERRED_USERNAME"
  auto_register        = false
}`, resourceName, frame.UniqueResourcesID, frame.OrgID, name, secret)
	}, idp_utils.ClientSecretVar)
}
