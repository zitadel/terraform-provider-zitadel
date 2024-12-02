package default_login_policy_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/admin"

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
		checkRemoteProperty(*frame),
		helper.ZitadelGeneratedIdOnlyRegex,
		test_utils.CheckNothing,
		test_utils.ImportNothing,
	)
}

func checkRemoteProperty(frame test_utils.InstanceTestFrame) func(string) resource.TestCheckFunc {
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
