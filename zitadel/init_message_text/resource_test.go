package init_message_text_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/management"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper/test_utils"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/init_message_text"
)

func TestAccInitMessageText(t *testing.T) {
	frame := test_utils.NewOrgTestFrame(t, "zitadel_init_message_text")
	resourceExample, exampleAttributes := test_utils.ReadExample(t, test_utils.Resources, frame.ResourceType)
	exampleProperty := test_utils.AttributeValue(t, "title", exampleAttributes).AsString()
	exampleLanguage := test_utils.AttributeValue(t, init_message_text.LanguageVar, exampleAttributes).AsString()
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
		checkRemoteProperty(frame, exampleLanguage)("Zitadel - Initialize User"),
		nil,
	)
}

func checkRemoteProperty(frame *test_utils.OrgTestFrame, lang string) func(string) resource.TestCheckFunc {
	return func(expect string) resource.TestCheckFunc {
		return func(state *terraform.State) error {
			remoteResource, err := frame.GetCustomInitMessageText(frame, &management.GetCustomInitMessageTextRequest{Language: lang})
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
