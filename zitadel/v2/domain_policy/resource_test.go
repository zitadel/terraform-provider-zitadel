package domain_policy_test

import (
	"fmt"
	"testing"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/helper"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/management"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/helper/test_utils"
)

func TestAccDomainPolicy(t *testing.T) {
	resourceName := "zitadel_domain_policy"
	initialProperty := false
	updatedProperty := true
	frame, err := test_utils.NewOrgTestFrame(resourceName)
	if err != nil {
		t.Fatalf("setting up test context failed: %v", err)
	}
	otherFrame, err := frame.AnotherOrg("domain-policy-org-" + frame.UniqueResourcesID)
	if err != nil {
		t.Fatalf("setting up test context failed: %v", err)
	}
	test_utils.RunLifecyleTest[bool](
		t,
		otherFrame.BaseTestFrame,
		func(configProperty bool, _ string) string {
			return fmt.Sprintf(`
resource "%s" "%s" {
  org_id = "%s"
  user_login_must_be_domain                   = %t
  validate_org_domains                        = false
  smtp_sender_address_matches_instance_domain = false
}`, resourceName, otherFrame.UniqueResourcesID, otherFrame.OrgID, configProperty)
		},
		initialProperty, updatedProperty,
		"", "", "",
		false,
		checkRemoteProperty(*otherFrame),
		helper.ZitadelGeneratedIdOnlyRegex,
		checkRemoteProperty(*otherFrame)(initialProperty),
		nil,
	)
}

func checkRemoteProperty(frame test_utils.OrgTestFrame) func(bool) resource.TestCheckFunc {
	return func(expect bool) resource.TestCheckFunc {
		return func(state *terraform.State) error {
			resp, err := frame.GetDomainPolicy(frame, &management.GetDomainPolicyRequest{})
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
