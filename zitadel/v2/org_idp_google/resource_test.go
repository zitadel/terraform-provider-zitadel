package org_idp_google_test

import (
	"fmt"
	"testing"

	test_utils_org "github.com/zitadel/terraform-provider-zitadel/zitadel/v2/org_idp_utils/test_utils"
)

func TestAccZITADELOrgIdPGoogle(t *testing.T) {
	resourceName := "zitadel_org_idp_google"
	frame, err := test_utils_org.NewOrgTestFrame(resourceName)
	if err != nil {
		t.Fatalf("setting up test context failed: %v", err)
	}
	test_utils_org.RunBasicLifecyleTest(t, frame, func(name, secret string) string {
		return fmt.Sprintf(`
resource "%s" "%s" {
  org_id              = "%s"
  name                = "%s"
  client_id           = "aclientid"
  client_secret       = "%s"
  scopes              = ["two", "scopes"]
  is_linking_allowed  = false
  is_creation_allowed = true
  is_auto_creation    = false
  is_auto_update      = true
}`, resourceName, frame.UniqueResourcesID, frame.OrgID, name, secret)
	})
}
