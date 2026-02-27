package login_policy_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/management"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper/test_utils"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/idp_azure_ad/idp_azure_ad_test_dep"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/idp_google/idp_google_test_dep"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/login_policy"
)

func TestAccLoginPolicy(t *testing.T) {
	frame := test_utils.NewOrgTestFrame(t, "zitadel_login_policy")
	resourceExample, exampleAttributes := test_utils.ReadExample(t, test_utils.Resources, frame.ResourceType)
	exampleProperty := test_utils.AttributeValue(t, login_policy.DefaultRedirectURIVar, exampleAttributes).AsString()
	azureADDep, _ := idp_azure_ad_test_dep.Create(t, frame.BaseTestFrame, frame.Admin)
	googleDep, _ := idp_google_test_dep.Create(t, frame.BaseTestFrame, frame.Admin)
	test_utils.RunLifecyleTest(
		t,
		frame.BaseTestFrame,
		[]string{frame.AsOrgDefaultDependency, azureADDep, googleDep},
		test_utils.ReplaceAll(resourceExample, exampleProperty, ""),
		exampleProperty, "localhost:9090",
		"", "", "",
		false,
		checkRemoteProperty(frame),
		helper.ZitadelGeneratedIdOnlyRegex,
		checkRemoteProperty(frame)(""),
		test_utils.ImportOrgId(frame),
	)
}

func TestAccLoginPolicySecondFactorsUpdate(t *testing.T) {
	frame := test_utils.NewOrgTestFrame(t, "zitadel_login_policy")

	initialConfig := fmt.Sprintf(`
%s
%s
resource "zitadel_login_policy" "default" {
  org_id                        = data.zitadel_org.default.id
  user_login                    = true
  allow_register                = false
  allow_external_idp            = false
  force_mfa                     = false
  force_mfa_local_only          = false
  passwordless_type             = "PASSWORDLESS_TYPE_ALLOWED"
  hide_password_reset           = false
  password_check_lifetime       = "240h0m0s"
  external_login_check_lifetime = "240h0m0s"
  multi_factor_check_lifetime   = "24h0m0s"
  mfa_init_skip_lifetime        = "720h0m0s"
  second_factor_check_lifetime  = "24h0m0s"
  ignore_unknown_usernames      = false
  default_redirect_uri          = "localhost:8080"
  second_factors                = ["SECOND_FACTOR_TYPE_OTP"]
}
`, frame.ProviderSnippet, frame.AsOrgDefaultDependency)

	updatedConfig := fmt.Sprintf(`
%s
%s
resource "zitadel_login_policy" "default" {
  org_id                        = data.zitadel_org.default.id
  user_login                    = true
  allow_register                = false
  allow_external_idp            = false
  force_mfa                     = false
  force_mfa_local_only          = false
  passwordless_type             = "PASSWORDLESS_TYPE_ALLOWED"
  hide_password_reset           = false
  password_check_lifetime       = "240h0m0s"
  external_login_check_lifetime = "240h0m0s"
  multi_factor_check_lifetime   = "24h0m0s"
  mfa_init_skip_lifetime        = "720h0m0s"
  second_factor_check_lifetime  = "24h0m0s"
  ignore_unknown_usernames      = false
  default_redirect_uri          = "localhost:8080"
  second_factors                = ["SECOND_FACTOR_TYPE_OTP", "SECOND_FACTOR_TYPE_U2F"]
  multi_factors                 = ["MULTI_FACTOR_TYPE_U2F_WITH_VERIFICATION"]
}
`, frame.ProviderSnippet, frame.AsOrgDefaultDependency)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: frame.V6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: initialConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(frame.TerraformName, "second_factors.#", "1"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(frame.TerraformName, "second_factors.#", "2"),
					resource.TestCheckResourceAttr(frame.TerraformName, "multi_factors.#", "1"),
				),
			},
		},
	})
}

func checkRemoteProperty(frame *test_utils.OrgTestFrame) func(string) resource.TestCheckFunc {
	return func(expect string) resource.TestCheckFunc {
		return func(state *terraform.State) error {
			resp, err := frame.GetLoginPolicy(frame, &management.GetLoginPolicyRequest{})
			if err != nil {
				return fmt.Errorf("getting policy failed: %w", err)
			}
			actual := resp.GetPolicy().GetDefaultRedirectUri()
			if actual != expect {
				return fmt.Errorf("expected %s, but got %s", expect, actual)
			}
			return nil
		}
	}
}
