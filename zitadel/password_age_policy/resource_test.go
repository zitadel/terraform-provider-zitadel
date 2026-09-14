package password_age_policy_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zclconf/go-cty/cty/gocty"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/management"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper/test_utils"
)

func TestPasswordAgePolicy(t *testing.T) {
	frame := test_utils.NewOrgTestFrame(t, "zitadel_password_age_policy")
	resourceExample, exampleAttributes := test_utils.ReadExample(t, test_utils.Resources, frame.ResourceType)
	ctyVal := test_utils.AttributeValue(t, "max_age_days", exampleAttributes)

	var maxAgeDays uint64

	if err := gocty.FromCtyValue(ctyVal, &maxAgeDays); err != nil {
		t.Fatalf("could not parse max_age_days: %s", err)
	}

	updatedProperty := uint64(15)
	test_utils.RunLifecyleTest(
		t,
		frame.BaseTestFrame,
		[]string{frame.AsOrgDefaultDependency},
		test_utils.ReplaceAll(resourceExample, maxAgeDays, ""),
		maxAgeDays, updatedProperty,
		"", "", "",
		false,
		checkRemoteProperty(frame),
		helper.ZitadelGeneratedIdOnlyRegex,
		test_utils.CheckNothing,
		test_utils.ImportOrgId(frame),
	)
}

func checkRemoteProperty(frame *test_utils.OrgTestFrame) func(uint64) resource.TestCheckFunc {
	return func(expect uint64) resource.TestCheckFunc {
		return func(state *terraform.State) error {
			resp, err := frame.GetPasswordAgePolicy(frame, &management.GetPasswordAgePolicyRequest{})
			if err != nil {
				return fmt.Errorf("getting policy failed: %w", err)
			}
			actual := resp.GetPolicy().GetMaxAgeDays()
			if actual != expect {
				return fmt.Errorf("expected %d, but got %d", expect, actual)
			}
			return nil
		}
	}
}

// TestAccPasswordAgePolicyWithoutOrgID reproduces #436 for this resource: creating the policy without
// org_id must record the organization of the authenticated service account
// instead of leaving the resource out of state.
func TestAccPasswordAgePolicyWithoutOrgID(t *testing.T) {
	frame := test_utils.NewOrgTestFrame(t, "zitadel_password_age_policy")

	resetToDefault := func() {
		if _, err := frame.ResetPasswordAgePolicyToDefault(frame, &management.ResetPasswordAgePolicyToDefaultRequest{}); err != nil {
			t.Logf("resetting policy to default: %v", err)
		}
	}
	resetToDefault()
	t.Cleanup(resetToDefault)

	config := fmt.Sprintf(`
%s
resource "zitadel_password_age_policy" "default" {
  max_age_days     = 30
  expire_warn_days = 5
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
					checkRemoteProperty(frame)(30),
				),
			},
			{
				Config:   config,
				PlanOnly: true,
			},
		},
	})
}
