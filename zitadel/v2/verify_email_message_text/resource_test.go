package verify_email_message_text_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/management"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/helper/test_utils"
)

func TestAccVerifyEmailMessageText(t *testing.T) {
	resourceName := "zitadel_verify_email_message_text"
	initialProperty := "initialtitle"
	updatedProperty := "updatedtitle"
	language := "en"
	frame, err := test_utils.NewOrgTestFrame(resourceName)
	if err != nil {
		t.Fatalf("setting up test context failed: %v", err)
	}
	test_utils.RunLifecyleTest[string](
		t,
		frame.BaseTestFrame,
		func(configProperty, _ string) string {
			return fmt.Sprintf(`
resource "%s" "%s" {
  org_id      = "%s"
  language    = "%s"

  title       = "%s"
  pre_header  = "pre_header example"
  subject     = "subject example"
  greeting    = "greeting example"
  text        = "text example"
  button_text = "button_text example"
  footer_text = "footer_text example"
}`, resourceName, frame.UniqueResourcesID, frame.OrgID, language, configProperty)
		},
		initialProperty, updatedProperty,
		"", "", "",
		true,
		checkRemoteProperty(frame, language),
		regexp.MustCompile(`^\d{18}_en$`),
		// When deleted, the default should be returned
		checkRemoteProperty(frame, language)("Zitadel - Verify email"),
		nil,
	)
}

func checkRemoteProperty(frame *test_utils.OrgTestFrame, lang string) func(string) resource.TestCheckFunc {
	return func(expect string) resource.TestCheckFunc {
		return func(state *terraform.State) error {
			remoteResource, err := frame.GetCustomVerifyEmailMessageText(frame, &management.GetCustomVerifyEmailMessageTextRequest{Language: lang})
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
