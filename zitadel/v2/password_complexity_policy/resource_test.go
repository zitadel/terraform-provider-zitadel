package password_complexity_policy_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/management"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/helper"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/helper/test_utils"
)

func TestAccPasswordComplexityPolicy(t *testing.T) {
	frame := test_utils.NewOrgTestFrame(t, "zitadel_password_complexity_policy")
	resourceExample, exampleAttributes := test_utils.ReadExample(t, test_utils.Resources, frame.ResourceType)
	exampleProperty, err := strconv.ParseUint(test_utils.AttributeValue(t, "min_length", exampleAttributes).AsString(), 10, 64)
	if err != nil {
		t.Fatalf("could not parse example property: %v", err)
	}
	updatedProperty := uint64(10)
	test_utils.RunLifecyleTest(
		t,
		frame.BaseTestFrame,
		[]string{frame.AsOrgDefaultDependency},
		test_utils.ReplaceAll(resourceExample, exampleProperty, ""),
		exampleProperty, updatedProperty,
		"", "",
		false,
		checkRemoteProperty(*frame),
		helper.ZitadelGeneratedIdOnlyRegex,
		checkRemoteProperty(*frame)(exampleProperty),
		nil, nil, "", "",
	)
}

func checkRemoteProperty(frame test_utils.OrgTestFrame) func(uint64) resource.TestCheckFunc {
	return func(expect uint64) resource.TestCheckFunc {
		return func(state *terraform.State) error {
			resp, err := frame.GetPasswordComplexityPolicy(frame, &management.GetPasswordComplexityPolicyRequest{})
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
