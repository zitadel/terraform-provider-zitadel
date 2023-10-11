package verify_sms_otp_message_text_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/management"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/helper/test_utils"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/verify_sms_otp_message_text"
)

func TestAccVerifySMSOTPMessageText(t *testing.T) {
	resourceName := "zitadel_verify_sms_otp_message_text"
	frame := test_utils.NewOrgTestFrame(t, resourceName)
	resourceExample, exampleAttributes := test_utils.ReadExample(t, test_utils.Resources, frame.ResourceType)
	exampleProperty := test_utils.AttributeValue(t, "text", exampleAttributes).AsString()
	exampleLanguage := test_utils.AttributeValue(t, verify_sms_otp_message_text.LanguageVar, exampleAttributes).AsString()
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
		checkRemotePropertyNotEmpty(frame, exampleLanguage),
		nil,
	)
}

func checkRemoteProperty(frame *test_utils.OrgTestFrame, lang string) func(string) resource.TestCheckFunc {
	return func(expect string) resource.TestCheckFunc {
		return func(state *terraform.State) error {
			remoteResource, err := frame.GetCustomVerifySMSOTPMessageText(frame, &management.GetCustomVerifySMSOTPMessageTextRequest{Language: lang})
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

func checkRemotePropertyNotEmpty(frame *test_utils.OrgTestFrame, lang string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		remoteResource, err := frame.GetCustomVerifySMSOTPMessageText(frame, &management.GetCustomVerifySMSOTPMessageTextRequest{Language: lang})
		if err != nil {
			return err
		}
		if remoteResource.GetCustomText().GetText() == "" {
			return fmt.Errorf("expected text not empty, but got empty")
		}
		return nil
	}
}
