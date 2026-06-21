package default_verify_sms_otp_message_text_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/admin"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/default_verify_sms_otp_message_text"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper/test_utils"
)

func TestAccDefaultVerifySMSOTPMessageText(t *testing.T) {
	frame := test_utils.NewInstanceTestFrame(t, "zitadel_default_verify_sms_otp_message_text")
	resourceExample, exampleAttributes := test_utils.ReadExample(t, test_utils.Resources, frame.ResourceType)
	exampleProperty := test_utils.AttributeValue(t, "text", exampleAttributes).AsString()
	exampleLanguage := test_utils.AttributeValue(t, default_verify_sms_otp_message_text.LanguageVar, exampleAttributes).AsString()
	test_utils.RunLifecyleTest(
		t,
		frame.BaseTestFrame,
		nil,
		test_utils.ReplaceAll(resourceExample, exampleProperty, ""),
		exampleProperty, "updatedtext",
		"", "", "",
		true,
		checkRemoteProperty(frame, exampleLanguage),
		regexp.MustCompile(fmt.Sprintf(`^%s$`, exampleLanguage)),
		// When deleted, the default should be returned
		checkRemotePropertyNotEmpty(frame, exampleLanguage),
		nil,
	)
}

// TestAccDefaultVerifySMSOTPMessageText_DeprecatedFieldsConverge asserts that
// setting the deprecated fields the SMS OTP API does not persist does not
// produce a perpetual diff (regression test for #412). The plan-only step uses
// the default ExpectNonEmptyPlan=false, so a non-empty plan after apply fails.
func TestAccDefaultVerifySMSOTPMessageText_DeprecatedFieldsConverge(t *testing.T) {
	frame := test_utils.NewInstanceTestFrame(t, "zitadel_default_verify_sms_otp_message_text")
	config := fmt.Sprintf(`%s
resource "zitadel_default_verify_sms_otp_message_text" "default" {
  language    = "en"
  text        = "text example"
  greeting    = "Greeting"
  subject     = "Subject"
  title       = "Title"
  pre_header  = "Pre header"
  button_text = "Button text"
  footer_text = "Footer text"
}
`, frame.ProviderSnippet)
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: frame.V6ProviderFactories(),
		Steps: []resource.TestStep{
			{Config: config},
			{Config: config, PlanOnly: true},
		},
	})
}

func checkRemoteProperty(frame *test_utils.InstanceTestFrame, lang string) func(string) resource.TestCheckFunc {
	return func(expect string) resource.TestCheckFunc {
		return func(state *terraform.State) error {
			remoteResource, err := frame.GetCustomVerifySMSOTPMessageText(frame, &admin.GetCustomVerifySMSOTPMessageTextRequest{Language: lang})
			if err != nil {
				return err
			}
			actual := remoteResource.GetCustomText().GetText()
			if actual != expect {
				return fmt.Errorf("expected %s, but got %s", expect, actual)
			}
			return nil
		}
	}
}
func checkRemotePropertyNotEmpty(frame *test_utils.InstanceTestFrame, lang string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		remoteResource, err := frame.GetCustomVerifySMSOTPMessageText(frame, &admin.GetCustomVerifySMSOTPMessageTextRequest{Language: lang})
		if err != nil {
			return err
		}
		if remoteResource.GetCustomText().GetText() == "" {
			return fmt.Errorf("expected text not empty, but got empty")
		}
		return nil
	}
}
