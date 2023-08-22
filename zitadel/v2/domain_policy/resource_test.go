package domain_policy_test

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/management"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/domain_policy"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/helper"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/helper/test_utils"
)

func TestAccDomainPolicy(t *testing.T) {
	frame := test_utils.NewOrgTestFrame(t, "zitadel_domain_policy")
	otherFrame := frame.AnotherOrg(t, "domain-policy-org-"+frame.UniqueResourcesID)
	resourceExample, resourceAttributes := test_utils.ReadExample(t, test_utils.Resources, frame.ResourceType)
	exampleProperty := test_utils.AttributeValue(t, domain_policy.UserLoginMustBeDomainVar, resourceAttributes).True()
	test_utils.RunLifecyleTest(
		t,
		otherFrame.BaseTestFrame,
		[]string{otherFrame.AsOrgDefaultDependency},
		func(property bool, secret string) string {
			// only replace first bool for the smtp_sender_address_matches_instance_domain property
			return strings.Replace(resourceExample, strconv.FormatBool(exampleProperty), strconv.FormatBool(property), 1)
		},
		exampleProperty, !exampleProperty,
		"", "", "",
		false,
		checkRemoteProperty(*otherFrame),
		helper.ZitadelGeneratedIdOnlyRegex,
		checkRemoteProperty(*otherFrame)(false),
		test_utils.ImportOrgId(otherFrame),
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
