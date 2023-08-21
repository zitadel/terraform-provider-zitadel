package verify_phone_message_text_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/management"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/helper/test_utils"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/verify_phone_message_text"
)

func TestAccVerifyPhoneMessageText(t *testing.T) {
	frame := test_utils.NewOrgTestFrame(t, "zitadel_verify_phone_message_text")
	resourceExample, exampleAttributes := test_utils.ReadExample(t, test_utils.Resources, frame.ResourceType)
	exampleProperty := test_utils.AttributeValue(t, "title", exampleAttributes).AsString()
	exampleLanguage := test_utils.AttributeValue(t, verify_phone_message_text.LanguageVar, exampleAttributes).AsString()
	test_utils.RunLifecyleTest(
		t,
		frame.BaseTestFrame,
		[]string{frame.AsOrgDefaultDependency},
		test_utils.ReplaceAll(resourceExample, exampleProperty, ""),
		exampleProperty, "updatedtitle",
		"", "", "",
		true,
		checkRemoteProperty(frame, exampleLanguage),
		regexp.MustCompile(fmt.Sprintf(`^\d{18}_%s$`, exampleLanguage)),
		// When deleted, the default should be returned
		checkRemoteProperty(frame, exampleLanguage)("Zitadel - Verify phone"),
		nil,
	)
}

func checkRemoteProperty(frame *test_utils.OrgTestFrame, lang string) func(string) resource.TestCheckFunc {
	return func(expect string) resource.TestCheckFunc {
		return func(state *terraform.State) error {
			remoteResource, err := frame.GetCustomVerifyPhoneMessageText(frame, &management.GetCustomVerifyPhoneMessageTextRequest{Language: lang})
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
