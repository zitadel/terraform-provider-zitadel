package default_hosted_login_translation_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	settingsv2 "github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/settings/v2"
	"google.golang.org/protobuf/types/known/structpb"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/default_hosted_login_translation"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper/test_utils"
)

func TestAccDefaultHostedLoginTranslation(t *testing.T) {
	frame := test_utils.NewInstanceTestFrame(t, "zitadel_default_hosted_login_translation")
	resourceExample, exampleAttributes := test_utils.ReadExample(t, test_utils.Resources, frame.ResourceType)
	exampleProperty := "Welcome to ACME"
	exampleLanguage := test_utils.AttributeValue(t, default_hosted_login_translation.LanguageVar, exampleAttributes).AsString()
	test_utils.RunLifecyleTest(
		t,
		frame.BaseTestFrame,
		nil,
		test_utils.ReplaceAll(resourceExample, exampleProperty, ""),
		exampleProperty, "updatedtitle",
		"", "", "",
		false,
		checkRemoteProperty(frame, exampleLanguage),
		regexp.MustCompile(fmt.Sprintf(`^%s$`, exampleLanguage)),
		// ZITADEL has no API to remove translations, so nothing changes remotely on destroy
		test_utils.CheckNothing,
		test_utils.ImportResourceId(frame.BaseTestFrame),
	)
}

// TestAccDefaultHostedLoginTranslationReplacesRemote verifies that the configured
// translations replace everything previously set for the language, and that a
// differently formatted but equivalent JSON document does not show up as a diff.
func TestAccDefaultHostedLoginTranslationReplacesRemote(t *testing.T) {
	frame := test_utils.NewInstanceTestFrame(t, "zitadel_default_hosted_login_translation")
	exampleLanguage := "en"
	title := "title " + frame.UniqueResourcesID

	resourceConfig := fmt.Sprintf(`
%s
resource "zitadel_default_hosted_login_translation" "default" {
  language     = "%s"
  translations = <<-EOT
    {
      "loginname": {
        "title": "%s"
      }
    }
  EOT
}
`, frame.ProviderSnippet, exampleLanguage, title)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: frame.V6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					client, err := helper.GetSecuritySettingsClient(context.Background(), frame.ClientInfo)
					if err != nil {
						t.Fatalf("failed to get client: %v", err)
					}
					translations, err := structpb.NewStruct(map[string]interface{}{
						"loginname": map[string]interface{}{"title": "remote", "description": "remote"},
					})
					if err != nil {
						t.Fatalf("failed to build translations: %v", err)
					}
					if _, err := client.SetHostedLoginTranslation(context.Background(), &settingsv2.SetHostedLoginTranslationRequest{
						Level:        &settingsv2.SetHostedLoginTranslationRequest_Instance{Instance: true},
						Locale:       exampleLanguage,
						Translations: translations,
					}); err != nil {
						t.Fatalf("setting remote hosted login translation failed: %v", err)
					}
				},
				Config: resourceConfig,
				Check: resource.ComposeTestCheckFunc(
					test_utils.CheckAMinute(checkRemoteProperty(frame, exampleLanguage)(title)),
					test_utils.CheckAMinute(checkRemoteDescription(frame, exampleLanguage, "")),
				),
			},
			{
				Config:   resourceConfig,
				PlanOnly: true,
			},
		},
	})
}

func getRemoteLoginname(frame *test_utils.InstanceTestFrame, language string) (*structpb.Struct, error) {
	client, err := helper.GetSecuritySettingsClient(context.Background(), frame.ClientInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to get client: %w", err)
	}
	resp, err := client.GetHostedLoginTranslation(context.Background(), &settingsv2.GetHostedLoginTranslationRequest{
		Level:             &settingsv2.GetHostedLoginTranslationRequest_Instance{Instance: true},
		Locale:            language,
		IgnoreInheritance: true,
	})
	if err != nil {
		return nil, fmt.Errorf("getting hosted login translation failed: %w", err)
	}
	return resp.GetTranslations().GetFields()["loginname"].GetStructValue(), nil
}

func checkRemoteProperty(frame *test_utils.InstanceTestFrame, language string) func(string) resource.TestCheckFunc {
	return func(expect string) resource.TestCheckFunc {
		return func(state *terraform.State) error {
			loginname, err := getRemoteLoginname(frame, language)
			if err != nil {
				return err
			}
			actual := loginname.GetFields()["title"].GetStringValue()
			if actual != expect {
				return fmt.Errorf("expected %q, but got %q", expect, actual)
			}
			return nil
		}
	}
}

func checkRemoteDescription(frame *test_utils.InstanceTestFrame, language, expect string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		loginname, err := getRemoteLoginname(frame, language)
		if err != nil {
			return err
		}
		actual := loginname.GetFields()["description"].GetStringValue()
		if actual != expect {
			return fmt.Errorf("expected description %q, but got %q", expect, actual)
		}
		return nil
	}
}
