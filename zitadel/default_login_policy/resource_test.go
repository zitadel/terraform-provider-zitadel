package default_login_policy_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/admin"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/policy"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/default_login_policy"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper/test_utils"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/idp_azure_ad/idp_azure_ad_test_dep"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/idp_google/idp_google_test_dep"
)

func TestAccDefaultLoginPolicy(t *testing.T) {
	frame := test_utils.NewInstanceTestFrame(t, "zitadel_default_login_policy")
	resourceExample, exampleAttributes := test_utils.ReadExample(t, test_utils.Resources, frame.ResourceType)
	exampleProperty := test_utils.AttributeValue(t, default_login_policy.DefaultRedirectURIVar, exampleAttributes).AsString()
	azureADDep, _ := idp_azure_ad_test_dep.Create(t, frame.BaseTestFrame, frame)
	googleDep, _ := idp_google_test_dep.Create(t, frame.BaseTestFrame, frame)
	test_utils.RunLifecyleTest(
		t,
		frame.BaseTestFrame,
		[]string{azureADDep, googleDep},
		test_utils.ReplaceAll(resourceExample, exampleProperty, ""),
		exampleProperty, "localhost:9090",
		"", "", "",
		false,
		checkRemoteProperty(frame),
		helper.ZitadelGeneratedIdOnlyRegex,
		test_utils.CheckNothing,
		test_utils.ImportNothing,
	)
}

func TestAccDefaultLoginPolicyCreateZeroValues(t *testing.T) {
	frame := test_utils.NewInstanceTestFrame(t, "zitadel_default_login_policy")

	zeroValuesConfig := fmt.Sprintf(`
%s
resource "zitadel_default_login_policy" "default" {
  user_login                    = true
  allow_register                = true
  allow_external_idp            = true
  force_mfa                     = false
  force_mfa_local_only          = false
  passwordless_type             = "PASSWORDLESS_TYPE_ALLOWED"
  hide_password_reset           = "false"
  password_check_lifetime       = "240h0m0s"
  external_login_check_lifetime = "240h0m0s"
  multi_factor_check_lifetime   = "24h0m0s"
  mfa_init_skip_lifetime        = "720h0m0s"
  second_factor_check_lifetime  = "24h0m0s"
  ignore_unknown_usernames      = true
  default_redirect_uri          = "localhost:8080"
  second_factors                = []
  multi_factors                 = []
  allow_domain_discovery        = true
  disable_login_with_email      = true
  disable_login_with_phone      = true
}
`, frame.ProviderSnippet)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: frame.V6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					if _, err := frame.AddSecondFactorToLoginPolicy(frame, &admin.AddSecondFactorToLoginPolicyRequest{
						Type: policy.SecondFactorType_SECOND_FACTOR_TYPE_OTP,
					}); helper.IgnoreAlreadyExistsError(err) != nil {
						t.Fatalf("adding remote second factor failed: %v", err)
					}
					if _, err := frame.AddMultiFactorToLoginPolicy(frame, &admin.AddMultiFactorToLoginPolicyRequest{
						Type: policy.MultiFactorType_MULTI_FACTOR_TYPE_U2F_WITH_VERIFICATION,
					}); helper.IgnoreAlreadyExistsError(err) != nil {
						t.Fatalf("adding remote multi factor failed: %v", err)
					}
				},
				Config: zeroValuesConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(frame.TerraformName, "second_factors.#", "0"),
					resource.TestCheckResourceAttr(frame.TerraformName, "multi_factors.#", "0"),
					test_utils.CheckAMinute(checkRemoteSecondFactors(frame)(0)),
					test_utils.CheckAMinute(checkRemoteMultiFactors(frame)(0)),
				),
			},
		},
	})
}

func checkRemoteSecondFactors(frame *test_utils.InstanceTestFrame) func(int) resource.TestCheckFunc {
	return func(expect int) resource.TestCheckFunc {
		return func(state *terraform.State) error {
			resp, err := frame.ListLoginPolicySecondFactors(frame, &admin.ListLoginPolicySecondFactorsRequest{})
			if err != nil {
				return fmt.Errorf("listing second factors failed: %w", err)
			}
			actual := len(resp.GetResult())
			if actual != expect {
				return fmt.Errorf("expected %d second factors, but got %d: %v", expect, actual, resp.GetResult())
			}
			return nil
		}
	}
}

func checkRemoteMultiFactors(frame *test_utils.InstanceTestFrame) func(int) resource.TestCheckFunc {
	return func(expect int) resource.TestCheckFunc {
		return func(state *terraform.State) error {
			resp, err := frame.ListLoginPolicyMultiFactors(frame, &admin.ListLoginPolicyMultiFactorsRequest{})
			if err != nil {
				return fmt.Errorf("listing multi factors failed: %w", err)
			}
			actual := len(resp.GetResult())
			if actual != expect {
				return fmt.Errorf("expected %d multi factors, but got %d: %v", expect, actual, resp.GetResult())
			}
			return nil
		}
	}
}

func checkRemoteProperty(frame *test_utils.InstanceTestFrame) func(string) resource.TestCheckFunc {
	return func(expect string) resource.TestCheckFunc {
		return func(state *terraform.State) error {
			resp, err := frame.GetLoginPolicy(frame, &admin.GetLoginPolicyRequest{})
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
