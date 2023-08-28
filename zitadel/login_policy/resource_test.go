package login_policy_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/management"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/helper"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/helper/test_utils"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/idp_azure_ad/idp_azure_ad_test_dep"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/idp_google/idp_google_test_dep"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/login_policy"
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
		checkRemoteProperty(*frame),
		helper.ZitadelGeneratedIdOnlyRegex,
		checkRemoteProperty(*frame)(""),
		test_utils.ImportOrgId(frame),
	)
}

func checkRemoteProperty(frame test_utils.OrgTestFrame) func(string) resource.TestCheckFunc {
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
