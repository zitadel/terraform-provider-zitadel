package default_login_policy_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/admin"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/helper/test_utils"
)

func TestAccDefaultLoginPolicy(t *testing.T) {
	resourceName := "zitadel_default_login_policy"
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
  user_login                    = %t
  allow_register                = true
  allow_external_idp            = true
  force_mfa                     = false
  passwordless_type             = "PASSWORDLESS_TYPE_ALLOWED"
  hide_password_reset           = "false"
  password_check_lifetime       = "240h0m0s"
  external_login_check_lifetime = "240h0m0s"
  multi_factor_check_lifetime   = "24h0m0s"
  mfa_init_skip_lifetime        = "720h0m0s"
  second_factor_check_lifetime  = "24h0m0s"
  ignore_unknown_usernames      = true
  default_redirect_uri          = "localhost:8080"
  second_factors                = ["SECOND_FACTOR_TYPE_OTP", "SECOND_FACTOR_TYPE_U2F"]
  multi_factors                 = ["MULTI_FACTOR_TYPE_U2F_WITH_VERIFICATION"]
  allow_domain_discovery        = true
  disable_login_with_email      = true
  disable_login_with_phone      = true
}`, resourceName, frame.UniqueResourcesID, configProperty)
		},
		initialProperty, updatedProperty,
		"", "",
		checkRemoteProperty(*frame),
		test_utils.ZITADEL_GENERATED_ID_REGEX,
		test_utils.CheckNothing,
		nil, nil, "", "",
	)
}

func checkRemoteProperty(frame test_utils.InstanceTestFrame) func(bool) resource.TestCheckFunc {
	return func(expect bool) resource.TestCheckFunc {
		return func(state *terraform.State) error {
			resp, err := frame.GetLoginPolicy(frame, &admin.GetLoginPolicyRequest{})
			if err != nil {
				return fmt.Errorf("getting policy failed: %w", err)
			}
			actual := resp.GetPolicy().GetAllowUsernamePassword()
			if actual != expect {
				return fmt.Errorf("expected %t, but got %t", expect, actual)
			}
			return nil
		}
	}
}
