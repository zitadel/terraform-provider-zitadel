package notification_policy_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/management"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper/test_utils"
)

func TestAccNotificationPolicy(t *testing.T) {
	frame := test_utils.NewOrgTestFrame(t, "zitadel_notification_policy")
	resourceExample, _ := test_utils.ReadExample(t, test_utils.Resources, frame.ResourceType)
	exampleProperty := false
	initialProperty := true
	updatedProperty := false
	test_utils.RunLifecyleTest(
		t,
		frame.BaseTestFrame,
		[]string{frame.AsOrgDefaultDependency},
		test_utils.ReplaceAll(resourceExample, exampleProperty, ""),
		initialProperty, updatedProperty,
		"", "", "",
		false,
		checkRemoteProperty(frame),
		helper.ZitadelGeneratedIdOnlyRegex,
		checkRemoteProperty(frame)(true),
		test_utils.ImportOrgId(frame),
	)
}

func checkRemoteProperty(frame *test_utils.OrgTestFrame) func(bool) resource.TestCheckFunc {
	return func(expect bool) resource.TestCheckFunc {
		return func(state *terraform.State) error {
			resp, err := frame.GetNotificationPolicy(frame, &management.GetNotificationPolicyRequest{})
			if err != nil {
				return fmt.Errorf("getting policy failed: %w", err)
			}
			actual := resp.GetPolicy().GetPasswordChange()
			if actual != expect {
				return fmt.Errorf("expected %t, but got %t", expect, actual)
			}
			return nil
		}
	}
}

// TestAccNotificationPolicyWithoutOrgID reproduces #436 for this resource: creating the policy without
// org_id must record the organization of the authenticated service account
// instead of leaving the resource out of state.
func TestAccNotificationPolicyWithoutOrgID(t *testing.T) {
	frame := test_utils.NewOrgTestFrame(t, "zitadel_notification_policy")

	resetToDefault := func() {
		if _, err := frame.ResetNotificationPolicyToDefault(frame, &management.ResetNotificationPolicyToDefaultRequest{}); err != nil {
			t.Logf("resetting policy to default: %v", err)
		}
	}
	resetToDefault()
	t.Cleanup(resetToDefault)

	config := fmt.Sprintf(`
%s
resource "zitadel_notification_policy" "default" {
  password_change = false
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
					checkRemoteProperty(frame)(false),
				),
			},
			{
				Config:   config,
				PlanOnly: true,
			},
		},
	})
}
