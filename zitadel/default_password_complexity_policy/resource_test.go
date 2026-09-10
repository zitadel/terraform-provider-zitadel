package default_password_complexity_policy_test

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/admin"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/default_password_complexity_policy"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper/test_utils"
)

func TestAccDefaultPasswordComplexityPolicy(t *testing.T) {
	frame := test_utils.NewInstanceTestFrame(t, "zitadel_default_password_complexity_policy")
	resourceExample, exampleAttributes := test_utils.ReadExample(t, test_utils.Resources, frame.ResourceType)
	exampleProperty, err := strconv.ParseUint(test_utils.AttributeValue(t, default_password_complexity_policy.MinLengthVar, exampleAttributes).AsString(), 10, 64)
	if err != nil {
		t.Fatalf("could not parse example property: %v", err)
	}
	test_utils.RunLifecyleTest(
		t,
		frame.BaseTestFrame,
		nil,
		test_utils.ReplaceAll(resourceExample, exampleProperty, ""),
		exampleProperty, 10,
		"", "", "",
		false,
		checkRemoteProperty(frame),
		helper.ZitadelGeneratedIdOnlyRegex,
		test_utils.CheckNothing,
		test_utils.ImportNothing,
	)
}

func TestAccDefaultPasswordComplexityPolicyCreateZeroValues(t *testing.T) {
	frame := test_utils.NewInstanceTestFrame(t, "zitadel_default_password_complexity_policy")

	zeroValuesConfig := fmt.Sprintf(`
%s
resource "zitadel_default_password_complexity_policy" "default" {
  min_length    = 0
  has_uppercase = false
  has_lowercase = false
  has_number    = false
  has_symbol    = false
}
`, frame.ProviderSnippet)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: frame.V6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					if _, err := frame.UpdatePasswordComplexityPolicy(frame, &admin.UpdatePasswordComplexityPolicyRequest{
						MinLength:    8,
						HasUppercase: true,
						HasLowercase: true,
						HasNumber:    true,
						HasSymbol:    true,
					}); helper.IgnorePreconditionError(err) != nil {
						t.Fatalf("setting remote policy failed: %v", err)
					}
				},
				Config:      zeroValuesConfig,
				ExpectError: regexp.MustCompile(`Given minimum length is not allowed`),
			},
		},
	})
}

func checkRemoteProperty(frame *test_utils.InstanceTestFrame) func(uint64) resource.TestCheckFunc {
	return func(expect uint64) resource.TestCheckFunc {
		return func(state *terraform.State) error {
			resp, err := frame.GetPasswordComplexityPolicy(frame, &admin.GetPasswordComplexityPolicyRequest{})
			if err != nil {
				return fmt.Errorf("getting policy failed: %w", err)
			}
			actual := resp.GetPolicy().GetMinLength()
			if actual != expect {
				return fmt.Errorf("expected %d, but got %d", expect, actual)
			}
			return nil
		}
	}
}
