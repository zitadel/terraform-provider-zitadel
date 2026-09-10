package default_domain_policy_test

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/admin"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/default_domain_policy"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper/test_utils"
)

func TestAccDefaultDomainPolicy(t *testing.T) {
	frame := test_utils.NewInstanceTestFrame(t, "zitadel_default_domain_policy")
	resourceExample, resourceAttributes := test_utils.ReadExample(t, test_utils.Resources, frame.ResourceType)
	exampleProperty := test_utils.AttributeValue(t, default_domain_policy.UserLoginMustBeDomainVar, resourceAttributes).True()
	test_utils.RunLifecyleTest(
		t,
		frame.BaseTestFrame,
		nil,
		func(property bool, secret string) string {
			// only replace first bool for the smtp_sender_address_matches_instance_domain property
			return strings.Replace(resourceExample, strconv.FormatBool(exampleProperty), strconv.FormatBool(property), 1)
		},
		exampleProperty, !exampleProperty,
		"", "", "",
		false,
		checkRemoteProperty(frame),
		helper.ZitadelGeneratedIdOnlyRegex,
		test_utils.CheckNothing,
		test_utils.ImportNothing,
	)
}

func TestAccDefaultDomainPolicyCreateZeroValues(t *testing.T) {
	frame := test_utils.NewInstanceTestFrame(t, "zitadel_default_domain_policy")

	zeroValuesConfig := fmt.Sprintf(`
%s
resource "zitadel_default_domain_policy" "default" {
  user_login_must_be_domain                   = false
  validate_org_domains                        = false
  smtp_sender_address_matches_instance_domain = false
}
`, frame.ProviderSnippet)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: frame.V6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					if _, err := frame.UpdateDomainPolicy(frame, &admin.UpdateDomainPolicyRequest{
						UserLoginMustBeDomain:                  true,
						ValidateOrgDomains:                     false,
						SmtpSenderAddressMatchesInstanceDomain: false,
					}); helper.IgnorePreconditionError(err) != nil {
						t.Fatalf("setting remote policy failed: %v", err)
					}
				},
				Config: zeroValuesConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(frame.TerraformName, "user_login_must_be_domain", "false"),
					test_utils.CheckAMinute(checkRemoteProperty(frame)(false)),
				),
			},
		},
	})
}

func checkRemoteProperty(frame *test_utils.InstanceTestFrame) func(bool) resource.TestCheckFunc {
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
