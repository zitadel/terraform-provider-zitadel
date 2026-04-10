package sms_provider_twilio_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/admin"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper/test_utils"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/sms_provider_twilio"
)

func TestAccSMSProviderTwilio(t *testing.T) {
	frame := test_utils.NewInstanceTestFrame(t, "zitadel_sms_provider_twilio")
	resourceExample, exampleAttributes := test_utils.ReadExample(t, test_utils.Resources, frame.ResourceType)
	exampleProperty := test_utils.AttributeValue(t, sms_provider_twilio.SenderNumberVar, exampleAttributes).AsString()
	exampleSecret := test_utils.AttributeValue(t, sms_provider_twilio.TokenVar, exampleAttributes).AsString()
	test_utils.RunLifecyleTest(
		t,
		frame.BaseTestFrame,
		nil,
		test_utils.ReplaceAll(resourceExample, exampleProperty, exampleSecret),
		exampleProperty, "987654321",
		sms_provider_twilio.TokenVar, exampleSecret, "updatedSecret",
		false,
		checkRemoteProperty(frame),
		helper.ZitadelGeneratedIdOnlyRegex,
		test_utils.CheckNothing,
		test_utils.ChainImportStateIdFuncs(
			test_utils.ImportResourceId(frame.BaseTestFrame),
			test_utils.ImportStateAttribute(frame.BaseTestFrame, sms_provider_twilio.TokenVar),
		),
	)
}

func TestAccSMSProviderTwilioActivation(t *testing.T) {
	frame := test_utils.NewInstanceTestFrame(t, "zitadel_sms_provider_twilio")

	initialConfig := fmt.Sprintf(`
%s
resource "zitadel_sms_provider_twilio" "default" {
  sid           = "test_sid"
  token         = "test_token"
  sender_number = "123456789"
  set_active    = false
}
`, frame.ProviderSnippet)

	activatedConfig := fmt.Sprintf(`
%s
resource "zitadel_sms_provider_twilio" "default" {
  sid           = "test_sid"
  token         = "test_token"
  sender_number = "123456789"
  set_active    = true
}
`, frame.ProviderSnippet)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: frame.V6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: initialConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(frame.TerraformName, "set_active", "false"),
				),
			},
			{
				Config: activatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(frame.TerraformName, "set_active", "true"),
				),
			},
		},
	})
}

func TestAccSMSProviderTwilioVerifyServiceSid(t *testing.T) {
	frame := test_utils.NewInstanceTestFrame(t, "zitadel_sms_provider_twilio")
	resourceConfig := fmt.Sprintf(`
%s
resource "zitadel_sms_provider_twilio" "default" {
  sid                = "test_sid"
  token              = "test_token"
  sender_number      = "123456789"
  verify_service_sid = "VA1234567890abcdef1234567890abcdef"
}
`, frame.ProviderSnippet)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: frame.V6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: resourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(frame.TerraformName, "verify_service_sid", "VA1234567890abcdef1234567890abcdef"),
				),
			},
		},
	})
}

func TestAccSMSProviderTwilioDescription(t *testing.T) {
	frame := test_utils.NewInstanceTestFrame(t, "zitadel_sms_provider_twilio")
	resourceConfig := fmt.Sprintf(`
%s
resource "zitadel_sms_provider_twilio" "default" {
  sid           = "test_sid"
  token         = "test_token"
  sender_number = "123456789"
  description   = "My Twilio provider"
}
`, frame.ProviderSnippet)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: frame.V6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: resourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(frame.TerraformName, "description", "My Twilio provider"),
				),
			},
		},
	})
}

func TestAccSMSProviderTwilioFieldUpdate(t *testing.T) {
	frame := test_utils.NewInstanceTestFrame(t, "zitadel_sms_provider_twilio")
	initialConfig := fmt.Sprintf(`
%s
resource "zitadel_sms_provider_twilio" "default" {
  sid                = "test_sid"
  token              = "test_token"
  sender_number      = "123456789"
  verify_service_sid = "VA1234567890abcdef1234567890abcdef"
  description        = "Initial description"
}
`, frame.ProviderSnippet)

	updatedConfig := fmt.Sprintf(`
%s
resource "zitadel_sms_provider_twilio" "default" {
  sid                = "test_sid"
  token              = "test_token"
  sender_number      = "123456789"
  verify_service_sid = "VA0987654321fedcba0987654321fedcba"
  description        = "Updated description"
}
`, frame.ProviderSnippet)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: frame.V6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: initialConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(frame.TerraformName, "verify_service_sid", "VA1234567890abcdef1234567890abcdef"),
					resource.TestCheckResourceAttr(frame.TerraformName, "description", "Initial description"),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(frame.TerraformName, "verify_service_sid", "VA0987654321fedcba0987654321fedcba"),
					resource.TestCheckResourceAttr(frame.TerraformName, "description", "Updated description"),
				),
			},
		},
	})
}

func checkRemoteProperty(frame *test_utils.InstanceTestFrame) func(string) resource.TestCheckFunc {
	return func(expect string) resource.TestCheckFunc {
		return func(state *terraform.State) error {
			resp, err := frame.GetSMSProvider(frame, &admin.GetSMSProviderRequest{Id: frame.State(state).ID})
			if err != nil {
				return fmt.Errorf("getting sms provider failed: %w", err)
			}
			actual := resp.GetConfig().GetTwilio().GetSenderNumber()
			if actual != expect {
				return fmt.Errorf("expected %s, but got %s", expect, actual)
			}
			return nil
		}
	}
}
