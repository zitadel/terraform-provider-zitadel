package default_verify_email_otp_message_text_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/admin"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/default_verify_email_otp_message_text"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/helper/test_utils"
)

func TestAccDefaultVerifyEmailOTPMessageText(t *testing.T) {
	frame := test_utils.NewInstanceTestFrame(t, "zitadel_default_verify_email_otp_message_text")
	resourceExample, exampleAttributes := test_utils.ReadExample(t, test_utils.Resources, frame.ResourceType)
	exampleProperty := test_utils.AttributeValue(t, "title", exampleAttributes).AsString()
	exampleLanguage := test_utils.AttributeValue(t, default_verify_email_otp_message_text.LanguageVar, exampleAttributes).AsString()
	test_utils.RunLifecyleTest(
		t,
		frame.BaseTestFrame,
		nil,
		test_utils.ReplaceAll(resourceExample, exampleProperty, ""),
		exampleProperty, "updatedtitle",
		"", "", "",
		true,
		checkRemoteProperty(frame, exampleLanguage),
		regexp.MustCompile(fmt.Sprintf(`^%s$`, exampleLanguage)),
		// When deleted, the default should be returned
		checkRemoteProperty(frame, exampleLanguage)("Verify One-Time Password"),
		nil,
	)
}

func checkRemoteProperty(frame *test_utils.InstanceTestFrame, lang string) func(string) resource.TestCheckFunc {
	return func(expect string) resource.TestCheckFunc {
		return func(state *terraform.State) error {
			remoteResource, err := frame.GetCustomVerifyEmailOTPMessageText(frame, &admin.GetCustomVerifyEmailOTPMessageTextRequest{Language: lang})
			if err != nil {
				return err
			}
			actual := remoteResource.GetCustomText().GetTitle()
			if actual != expect {
				return fmt.Errorf("expected %s, but got %s", expect, actual)
			}
			return nil
		}
	}
}
