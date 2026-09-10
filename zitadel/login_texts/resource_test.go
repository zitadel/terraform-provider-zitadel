package login_texts_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/management"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper/test_utils"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/login_texts"
)

func TestAccLoginTexts(t *testing.T) {
	frame := test_utils.NewOrgTestFrame(t, "zitadel_login_texts")
	resourceExample, exampleAttributes := test_utils.ReadExample(t, test_utils.Resources, frame.ResourceType)
	exampleProperty := "example"
	exampleLanguage := test_utils.AttributeValue(t, login_texts.LanguageVar, exampleAttributes).AsString()
	test_utils.RunLifecyleTest(
		t,
		frame.BaseTestFrame,
		[]string{frame.AsOrgDefaultDependency},
		test_utils.ReplaceAll(resourceExample, exampleProperty, ""),
		exampleProperty, "updatedtext",
		"", "", "",
		true,
		checkRemoteProperty(frame, exampleLanguage),
		regexp.MustCompile(fmt.Sprintf(`^\d{18}_%s$`, exampleLanguage)),
		// When deleted, the default should be returned
		checkRemoteProperty(frame, exampleLanguage)(""),
		nil,
	)
}

// TestAccLoginTextsPartialConfig verifies that a config which only sets some of
// the screen-text blocks converges. ZITADEL returns all the other blocks with
// empty fields, which must not show up as a diff, while a block that is
// explicitly configured as {} must be kept.
func TestAccLoginTextsPartialConfig(t *testing.T) {
	frame := test_utils.NewOrgTestFrame(t, "zitadel_login_texts")
	exampleLanguage := "en"

	resourceConfig := fmt.Sprintf(`
%s
%s
resource "zitadel_login_texts" "default" {
  org_id   = data.zitadel_org.default.id
  language = "%s"

  email_verification_done_text = {
    title       = "partial"
    description = "Your email has been verified."
  }

  login_text = {
    title = "Welcome"
  }

  password_change_done_text = {}
}
`, frame.ProviderSnippet, frame.AsOrgDefaultDependency, exampleLanguage)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: frame.V6ProviderFactories(),
		CheckDestroy:             test_utils.CheckAMinute(checkRemoteProperty(frame, exampleLanguage)("")),
		Steps: []resource.TestStep{
			{
				Config: resourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(frame.TerraformName, "email_verification_done_text.title", "partial"),
					resource.TestCheckResourceAttr(frame.TerraformName, "login_text.title", "Welcome"),
					resource.TestCheckNoResourceAttr(frame.TerraformName, "logout_text.title"),
					test_utils.CheckAMinute(checkRemoteProperty(frame, exampleLanguage)("partial")),
				),
			},
			{
				Config:   resourceConfig,
				PlanOnly: true,
			},
		},
	})
}

func checkRemoteProperty(frame *test_utils.OrgTestFrame, lang string) func(string) resource.TestCheckFunc {
	return func(expect string) resource.TestCheckFunc {
		return func(state *terraform.State) error {
			remoteResource, err := frame.GetCustomLoginTexts(frame, &management.GetCustomLoginTextsRequest{Language: lang})
			if err != nil {
				return err
			}
			actual := remoteResource.GetCustomText().GetEmailVerificationDoneText().GetTitle()
			if actual != expect {
				return fmt.Errorf("expected %s, but got %s", expect, actual)
			}
			return nil
		}
	}
}
