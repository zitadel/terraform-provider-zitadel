package default_privacy_policy_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/admin"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/helper/test_utils"
)

func TestAccDefaultPrivacyPolicy(t *testing.T) {
	resourceName := "zitadel_default_privacy_policy"
	updatedProperty := "https://zitadel.com"
	initialProperty := "https://google.com"
	frame, err := test_utils.NewInstanceTestFrame(resourceName)
	if err != nil {
		t.Fatalf("setting up test context failed: %v", err)
	}
	test_utils.RunLifecyleTest[string](
		t,
		frame.BaseTestFrame,
		func(configProperty, _ string) string {
			return fmt.Sprintf(`
resource "%s" "%s" {
  tos_link     = "https://google.com"
  privacy_link = "https://google.com"
  help_link    = "%s"
}`, resourceName, frame.UniqueResourcesID, configProperty)
		},
		initialProperty, updatedProperty,
		"", "",
		checkRemoteProperty(*frame),
		test_utils.ZITADEL_GENERATED_ID_REGEX,
		test_utils.CheckNothing,
		nil, nil, "", "",
	)
}

func checkRemoteProperty(frame test_utils.InstanceTestFrame) func(string) resource.TestCheckFunc {
	return func(expect string) resource.TestCheckFunc {
		return func(state *terraform.State) error {
			resp, err := frame.GetPrivacyPolicy(frame, &admin.GetPrivacyPolicyRequest{})
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
