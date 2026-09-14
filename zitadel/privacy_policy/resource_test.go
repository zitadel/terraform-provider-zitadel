package privacy_policy_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/management"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper/test_utils"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/privacy_policy"
)

func TestAccPrivacyPolicy(t *testing.T) {
	frame := test_utils.NewOrgTestFrame(t, "zitadel_privacy_policy")
	resourceExample, exampleAttributes := test_utils.ReadExample(t, test_utils.Resources, frame.ResourceType)
	exampleProperty := test_utils.AttributeValue(t, privacy_policy.HelpLinkVar, exampleAttributes).AsString()
	test_utils.RunLifecyleTest(
		t,
		frame.BaseTestFrame,
		[]string{frame.AsOrgDefaultDependency},
		test_utils.ReplaceAll(resourceExample, exampleProperty, ""),
		exampleProperty, "http://example.com/acctest",
		"", "", "",
		false,
		checkRemoteProperty(frame),
		helper.ZitadelGeneratedIdOnlyRegex,
		checkRemoteProperty(frame)(""),
		test_utils.ImportOrgId(frame),
	)
}

func checkRemoteProperty(frame *test_utils.OrgTestFrame) func(string) resource.TestCheckFunc {
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

// TestAccPrivacyPolicyWithoutOrgID reproduces #436 for this resource: creating the policy without
// org_id must record the organization of the authenticated service account
// instead of leaving the resource out of state.
func TestAccPrivacyPolicyWithoutOrgID(t *testing.T) {
	frame := test_utils.NewOrgTestFrame(t, "zitadel_privacy_policy")

	resetToDefault := func() {
		if _, err := frame.ResetPrivacyPolicyToDefault(frame, &management.ResetPrivacyPolicyToDefaultRequest{}); err != nil {
			t.Logf("resetting policy to default: %v", err)
		}
	}
	resetToDefault()
	t.Cleanup(resetToDefault)

	config := fmt.Sprintf(`
%s
resource "zitadel_privacy_policy" "default" {
  help_link = "https://example.com/help"
}
`, frame.ProviderSnippet)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: frame.V6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(frame.TerraformName, "id", frame.OrgID),
					resource.TestCheckResourceAttr(frame.TerraformName, "org_id", frame.OrgID),
					checkRemoteProperty(frame)("https://example.com/help"),
				),
			},
			{
				Config:   config,
				PlanOnly: true,
			},
		},
	})
}
