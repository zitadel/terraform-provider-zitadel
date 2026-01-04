package default_lockout_policy_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/admin"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/default_lockout_policy"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper/test_utils"
)

func TestAccDefaultLockoutPolicy(t *testing.T) {
	frame := test_utils.NewInstanceTestFrame(t, "zitadel_default_lockout_policy")
	resourceExample, exampleAttributes := test_utils.ReadExample(t, test_utils.Resources, frame.ResourceType)
	exampleProperty, err := strconv.ParseUint(test_utils.AttributeValue(t, default_lockout_policy.MaxPasswordAttemptsVar, exampleAttributes).AsString(), 10, 64)
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

func checkRemoteProperty(frame *test_utils.InstanceTestFrame) func(uint64) resource.TestCheckFunc {
	return func(expect uint64) resource.TestCheckFunc {
		return func(state *terraform.State) error {
			resp, err := frame.GetLockoutPolicy(frame, &admin.GetLockoutPolicyRequest{})
			if err != nil {
				return fmt.Errorf("getting policy failed: %w", err)
			}
			actual := resp.GetPolicy().GetMaxPasswordAttempts()
			if actual != expect {
				return fmt.Errorf("expected %d, but got %d", expect, actual)
			}
			return nil
		}
	}
}
