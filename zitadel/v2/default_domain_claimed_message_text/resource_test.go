package default_domain_claimed_message_text_test

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/admin"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/helper/test_utils"
)

func TestAccDefaultDomainClaimedMessageText(t *testing.T) {
	resourceName := "zitadel_default_domain_claimed_message_text"
	language := "en"
	frame, err := test_utils.NewInstanceTestFrame(resourceName)
	if err != nil {
		t.Fatalf("setting up test context failed: %v", err)
	}
	resourceExample, exampleAttributes := frame.ReadExample(t, test_utils.Resources, frame.ResourceType)
	exampleProperty := test_utils.AttributeValue(t, "title", exampleAttributes).AsString()
	updatedProperty := "updatedtitle"
	test_utils.RunLifecyleTest[string](
		t,
		frame.BaseTestFrame,
		func(configProperty, _ string) string {
			return strings.Replace(resourceExample, exampleProperty, configProperty, 1)
		},
		exampleProperty, updatedProperty,
		"", "",
		true,
		checkRemoteProperty(frame, language),
		regexp.MustCompile(`^en$`),
		// When deleted, the default should be returned
		checkRemoteProperty(frame, language)("ZITADEL - Domain has been claimed"),
		nil, nil, "", "",
	)
}

func checkRemoteProperty(frame *test_utils.InstanceTestFrame, lang string) func(string) resource.TestCheckFunc {
	return func(expect string) resource.TestCheckFunc {
		return func(state *terraform.State) error {
			remoteResource, err := frame.GetCustomDomainClaimedMessageText(frame, &admin.GetCustomDomainClaimedMessageTextRequest{Language: lang})
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
