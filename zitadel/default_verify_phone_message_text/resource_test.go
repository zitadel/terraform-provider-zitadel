package default_verify_phone_message_text_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/admin"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/default_verify_phone_message_text"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/helper/test_utils"
)

func TestAccDefaultVerifyPhoneMessageText(t *testing.T) {
	frame := test_utils.NewInstanceTestFrame(t, "zitadel_default_verify_phone_message_text")
	resourceExample, exampleAttributes := test_utils.ReadExample(t, test_utils.Resources, frame.ResourceType)
	exampleProperty := test_utils.AttributeValue(t, "text", exampleAttributes).AsString()
	exampleLanguage := test_utils.AttributeValue(t, default_verify_phone_message_text.LanguageVar, exampleAttributes).AsString()
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

func checkRemoteProperty(frame *test_utils.InstanceTestFrame, lang string) func(string) resource.TestCheckFunc {
	return func(expect string) resource.TestCheckFunc {
		return func(state *terraform.State) error {
			remoteResource, err := frame.GetCustomVerifyPhoneMessageText(frame, &admin.GetCustomVerifyPhoneMessageTextRequest{Language: lang})
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
		remoteResource, err := frame.GetCustomVerifyPhoneMessageText(frame, &admin.GetCustomVerifyPhoneMessageTextRequest{Language: lang})
		if err != nil {
			return err
		}
		if remoteResource.GetCustomText().GetText() == "" {
			return fmt.Errorf("expected text not empty, but got empty")
		}
		return nil
	}
}
