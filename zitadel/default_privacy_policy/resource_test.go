package default_privacy_policy_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/admin"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/default_privacy_policy"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper/test_utils"
)

func TestAccDefaultPrivacyPolicy(t *testing.T) {
	frame := test_utils.NewInstanceTestFrame(t, "zitadel_default_privacy_policy")
	resourceExample, exampleAttributes := test_utils.ReadExample(t, test_utils.Resources, frame.ResourceType)
	exampleProperty := test_utils.AttributeValue(t, default_privacy_policy.HelpLinkVar, exampleAttributes).AsString()
	test_utils.RunLifecyleTest(
		t,
		frame.BaseTestFrame,
		nil,
		test_utils.ReplaceAll(resourceExample, exampleProperty, ""),
		exampleProperty, "http://example.com/acctest",
		"", "", "",
		false,
		checkRemoteProperty(frame),
		helper.ZitadelGeneratedIdOnlyRegex,
		test_utils.CheckNothing,
		test_utils.ImportNothing,
	)
}

func TestAccDefaultPrivacyPolicyCreateZeroValues(t *testing.T) {
	frame := test_utils.NewInstanceTestFrame(t, "zitadel_default_privacy_policy")

	zeroValuesConfig := fmt.Sprintf(`
%s
resource "zitadel_default_privacy_policy" "default" {
  tos_link      = ""
  privacy_link  = ""
  help_link     = ""
  support_email = ""
}
`, frame.ProviderSnippet)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: frame.V6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					if _, err := frame.UpdatePrivacyPolicy(frame, &admin.UpdatePrivacyPolicyRequest{
						TosLink:      "https://example.com/tos",
						PrivacyLink:  "https://example.com/privacy",
						HelpLink:     "https://example.com/help",
						SupportEmail: "support@example.com",
					}); helper.IgnorePreconditionError(err) != nil {
						t.Fatalf("setting remote policy failed: %v", err)
					}
				},
				Config: zeroValuesConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(frame.TerraformName, "help_link", ""),
					test_utils.CheckAMinute(checkRemoteProperty(frame)("")),
				),
			},
		},
	})
}

func checkRemoteProperty(frame *test_utils.InstanceTestFrame) func(string) resource.TestCheckFunc {
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
