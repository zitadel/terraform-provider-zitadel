package privacy_policy_test

import (
	"fmt"
	"testing"

	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/management"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/helper/test_utils"
)

func TestAccPrivacyPolicy(t *testing.T) {
	resourceName := "zitadel_privacy_policy"
	updatedProperty := "https://example.com"
	initialProperty := "https://httpbin.org"
	frame, err := test_utils.NewOrgTestFrame(resourceName)
	if err != nil {
		t.Fatalf("setting up test context failed: %v", err)
	}
	test_utils.RunLifecyleTest[string](
		t,
		frame.BaseTestFrame,
		func(configProperty, _ string) string {
			return fmt.Sprintf(`
resource "%s" "%s" {
  org_id       = "%s"
  tos_link     = "https://google.com"
  privacy_link = "https://google.com"
  help_link    = "%s"
}`, resourceName, frame.UniqueResourcesID, frame.OrgID, configProperty)
		},
		initialProperty, updatedProperty,
		"", "",
		checkRemoteProperty(*frame),
		test_utils.ZITADEL_GENERATED_ID_REGEX,
		checkRemoteProperty(*frame)(""),
		nil, nil, "", "",
	)
}

func checkRemoteProperty(frame test_utils.OrgTestFrame) func(string) resource.TestCheckFunc {
	return func(expect string) resource.TestCheckFunc {
		return func(state *terraform.State) error {
			resp, err := frame.GetPrivacyPolicy(frame, &management.GetPrivacyPolicyRequest{})
			if err != nil {
				return fmt.Errorf("getting policy failed: %w", err)
			}
			actual := resp.GetPolicy().GetHelpLink()
			if actual != expect {
				return fmt.Errorf("expected %s, but got %s", expect, actual)
			}
			return nil
		}
	}
}
