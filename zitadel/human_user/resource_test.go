package human_user_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/management"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper/test_utils"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/human_user"
)

func TestAccHumanUser(t *testing.T) {
	frame := test_utils.NewOrgTestFrame(t, "zitadel_human_user")
	resourceExample, exampleAttributes := test_utils.ReadExample(t, test_utils.Resources, frame.ResourceType)
	exampleUsername := test_utils.AttributeValue(t, human_user.UserNameVar, exampleAttributes).AsString()
	resourceExample = strings.Replace(resourceExample, exampleUsername, frame.UniqueResourcesID, 1)
	exampleProperty := test_utils.AttributeValue(t, human_user.DisplayNameVar, exampleAttributes).AsString()
	updatedProperty := "updatedproperty"
	test_utils.RunLifecyleTest(
		t,
		frame.BaseTestFrame,
		[]string{frame.AsOrgDefaultDependency},
		test_utils.ReplaceAll(resourceExample, exampleProperty, ""),
		exampleProperty, updatedProperty,
		"", "", "",
		false,
		checkRemoteProperty(frame),
		helper.ZitadelGeneratedIdOnlyRegex,
		test_utils.CheckIsNotFoundFromPropertyCheck(checkRemoteProperty(frame), updatedProperty),
		test_utils.ChainImportStateIdFuncs(
			test_utils.ImportResourceId(frame.BaseTestFrame),
			test_utils.ImportOrgId(frame),
			test_utils.ImportStateAttribute(frame.BaseTestFrame, human_user.InitialPasswordVar),
		),
	)
}

func TestAccHumanUserEmailVerifiedDrift(t *testing.T) {
	frame := test_utils.NewOrgTestFrame(t, "zitadel_human_user")

	configWithVerified := fmt.Sprintf(`%s
resource "%s" "default" {
	org_id            = "%s"
	user_name         = "%s"
	first_name        = "Test"
	last_name         = "User"
	email             = "test@example.com"
	is_email_verified = true
	initial_password  = "Password1!"
}
`, frame.ProviderSnippet, frame.ResourceType, frame.OrgID, frame.UniqueResourcesID)

	configWithoutVerified := fmt.Sprintf(`%s
resource "%s" "default" {
	org_id           = "%s"
	user_name        = "%s"
	first_name       = "Test"
	last_name        = "User"
	email            = "test@example.com"
	initial_password = "Password1!"
}
`, frame.ProviderSnippet, frame.ResourceType, frame.OrgID, frame.UniqueResourcesID)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: frame.V6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: configWithVerified,
			},
			{
				Config: configWithoutVerified,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(frame.TerraformName, "is_email_verified", "true"),
				),
			},
			{
				Config:             configWithoutVerified,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

func TestAccHumanUserPhoneVerifiedDrift(t *testing.T) {
	frame := test_utils.NewOrgTestFrame(t, "zitadel_human_user")

	configWithVerified := fmt.Sprintf(`%s
resource "%s" "default" {
	org_id            = "%s"
	user_name         = "%s"
	first_name        = "Test"
	last_name         = "User"
	email             = "test@example.com"
	phone             = "+1234567890"
	is_phone_verified = true
	initial_password  = "Password1!"
}
`, frame.ProviderSnippet, frame.ResourceType, frame.OrgID, frame.UniqueResourcesID)

	configWithoutVerified := fmt.Sprintf(`%s
resource "%s" "default" {
	org_id           = "%s"
	user_name        = "%s"
	first_name       = "Test"
	last_name        = "User"
	email            = "test@example.com"
	phone            = "+1234567890"
	initial_password = "Password1!"
}
`, frame.ProviderSnippet, frame.ResourceType, frame.OrgID, frame.UniqueResourcesID)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: frame.V6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: configWithVerified,
			},
			{
				Config: configWithoutVerified,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(frame.TerraformName, "is_phone_verified", "true"),
				),
			},
			{
				Config:             configWithoutVerified,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

func TestAccHumanUserProfileVariations(t *testing.T) {
	frame := test_utils.NewOrgTestFrame(t, "zitadel_human_user")

	initialConfig := fmt.Sprintf(`%s
resource "%s" "default" {
	org_id             = "%s"
	user_name          = "%s"
	first_name         = "Test"
	last_name          = "User"
	email              = "profile-test@example.com"
	is_email_verified  = true
	nick_name          = "testy"
	preferred_language = "en"
	gender             = "GENDER_FEMALE"
	initial_password   = "Password1!"
}
`, frame.ProviderSnippet, frame.ResourceType, frame.OrgID, frame.UniqueResourcesID)

	updatedConfig := fmt.Sprintf(`%s
resource "%s" "default" {
	org_id             = "%s"
	user_name          = "%s"
	first_name         = "Test"
	last_name          = "User"
	email              = "profile-test@example.com"
	is_email_verified  = true
	nick_name          = "updated-nick"
	preferred_language = "de"
	gender             = "GENDER_MALE"
	initial_password   = "Password1!"
}
`, frame.ProviderSnippet, frame.ResourceType, frame.OrgID, frame.UniqueResourcesID)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: frame.V6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: initialConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(frame.TerraformName, "nick_name", "testy"),
					resource.TestCheckResourceAttr(frame.TerraformName, "preferred_language", "en"),
					resource.TestCheckResourceAttr(frame.TerraformName, "gender", "GENDER_FEMALE"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(frame.TerraformName, "nick_name", "updated-nick"),
					resource.TestCheckResourceAttr(frame.TerraformName, "preferred_language", "de"),
					resource.TestCheckResourceAttr(frame.TerraformName, "gender", "GENDER_MALE"),
				),
			},
		},
	})
}

func checkRemoteProperty(frame *test_utils.OrgTestFrame) func(string) resource.TestCheckFunc {
	return func(expect string) resource.TestCheckFunc {
		return func(state *terraform.State) error {
			remoteResource, err := frame.GetUserByID(frame, &management.GetUserByIDRequest{Id: frame.State(state).ID})
			if err != nil {
				return err
			}
			actual := remoteResource.GetUser().GetHuman().GetProfile().GetDisplayName()
			if actual != expect {
				return fmt.Errorf("expected %s, but got %s", expect, actual)
			}
			return nil
		}
	}
}
