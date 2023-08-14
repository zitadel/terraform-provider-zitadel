package default_domain_policy_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/admin"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/helper"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/helper/test_utils"
)

func TestAccDefaultDomainPolicy(t *testing.T) {
	resourceName := "zitadel_default_domain_policy"
	initialProperty := true
	updatedProperty := false
	frame, err := test_utils.NewInstanceTestFrame(resourceName)
	if err != nil {
		t.Fatalf("setting up test context failed: %v", err)
	}
	test_utils.RunLifecyleTest[bool](
		t,
		frame.BaseTestFrame,
		func(configProperty bool, _ string) string {
			return fmt.Sprintf(`
resource "%s" "%s" {
  user_login_must_be_domain                   = %t
  validate_org_domains                        = false
  smtp_sender_address_matches_instance_domain = false
}`, resourceName, frame.UniqueResourcesID, configProperty)
		},
		initialProperty, updatedProperty,
		"", "", "",
		false,
		checkRemoteProperty(*frame),
		helper.ZitadelGeneratedIdOnlyRegex,
		test_utils.CheckNothing,
		test_utils.ImportNothing,
	)
}

func checkRemoteProperty(frame test_utils.InstanceTestFrame) func(bool) resource.TestCheckFunc {
	return func(expect bool) resource.TestCheckFunc {
		return func(state *terraform.State) error {
			resp, err := frame.GetDomainPolicy(frame, &admin.GetDomainPolicyRequest{})
			if err != nil {
				return fmt.Errorf("getting policy failed: %w", err)
			}
			actual := resp.GetPolicy().GetUserLoginMustBeDomain()
			if actual != expect {
				return fmt.Errorf("expected %t, but got %t", expect, actual)
			}
			return nil
		}
	}
}
