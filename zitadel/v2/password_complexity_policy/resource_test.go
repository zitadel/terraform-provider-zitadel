package password_complexity_policy_test

import (
	"fmt"
	"testing"

	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/management"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/helper/test_utils"
)

func TestAccPasswordComplexityPolicy(t *testing.T) {
	resourceName := "zitadel_password_complexity_policy"
	initialProperty := true
	updatedProperty := false
	frame, err := test_utils.NewOrgTestFrame(resourceName)
	if err != nil {
		t.Fatalf("setting up test context failed: %v", err)
	}
	test_utils.RunLifecyleTest[bool](
		t,
		frame.BaseTestFrame,
		func(configProperty bool, _ string) string {
			return fmt.Sprintf(`
resource "%s" "%s" {
  org_id = "%s"
  min_length    = "8"
  has_uppercase = true
  has_lowercase = true
  has_number    = true
  has_symbol    = %t
}`, resourceName, frame.UniqueResourcesID, frame.OrgID, configProperty)
		},
		initialProperty, updatedProperty,
		"", "",
		checkRemoteProperty(*frame),
		test_utils.ZITADEL_GENERATED_ID_REGEX,
		checkRemoteProperty(*frame)(initialProperty),
		nil, nil, "", "",
	)
}

func checkRemoteProperty(frame test_utils.OrgTestFrame) func(bool) resource.TestCheckFunc {
	return func(expect bool) resource.TestCheckFunc {
		return func(state *terraform.State) error {
			resp, err := frame.GetPasswordComplexityPolicy(frame, &management.GetPasswordComplexityPolicyRequest{})
			if err != nil {
				return fmt.Errorf("getting policy failed: %w", err)
			}
			actual := resp.GetPolicy().GetHasSymbol()
			if actual != expect {
				return fmt.Errorf("expected %t, but got %t", expect, actual)
			}
			return nil
		}
	}
}
