package email_provider_smtp_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/admin"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/settings"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/email_provider_smtp"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper/test_utils"
)

func TestAccEmailSMTPProvider(t *testing.T) {
	t.Skip("Skipping: flaky due to email provider state conflicts with email_provider_http on same instance")
	frame := test_utils.NewInstanceTestFrame(t, "zitadel_email_provider_smtp")
	resourceExample, exampleAttributes := test_utils.ReadExample(t, test_utils.Resources, frame.ResourceType)
	senderAddressProperty := test_utils.AttributeValue(t, email_provider_smtp.SenderAddressVar, exampleAttributes).AsString()
	resourceExample = strings.Replace(resourceExample, senderAddressProperty, fmt.Sprintf("zitadel@%s", frame.InstanceDomain), 1)

	exampleProperty := test_utils.AttributeValue(t, email_provider_smtp.SenderNameVar, exampleAttributes).AsString()
	updatedProperty := "updatedProperty"

	exampleSecret := test_utils.AttributeValue(t, email_provider_smtp.PasswordVar, exampleAttributes).AsString()
	updatedSecret := "updatedSecret"

	test_utils.RunLifecyleTest(
		t,
		frame.BaseTestFrame,
		nil,
		test_utils.ReplaceAll(resourceExample, exampleProperty, exampleSecret),
		exampleProperty, updatedProperty,
		email_provider_smtp.PasswordVar, exampleSecret, updatedSecret,
		false,
		checkRemoteProperty(frame),
		helper.ZitadelGeneratedIdOnlyRegex,
		test_utils.CheckIsNotFoundFromPropertyCheck(checkRemoteProperty(frame), ""),
		test_utils.ChainImportStateIdFuncs(
			test_utils.ImportResourceId(frame.BaseTestFrame),
		),
		// The password is write-only: it is never persisted to state, so it
		// cannot be sourced from state for the import ID, nor verified after
		// import. Its companion hash attribute is likewise absent after import.
		email_provider_smtp.PasswordVar, "password_hash", "set_active",
	)
}

func checkRemoteProperty(frame *test_utils.InstanceTestFrame) func(string) resource.TestCheckFunc {
	return func(expect string) resource.TestCheckFunc {
		return func(state *terraform.State) error {
			resp, err := frame.GetEmailProviderById(frame, &admin.GetEmailProviderByIdRequest{Id: frame.State(state).ID})
			if err != nil {
				return fmt.Errorf("getting email provider failed: %w", err)
			}
			actual := resp.GetConfig().GetSmtp().GetSenderName()
			if actual != expect {
				return fmt.Errorf("expected %s, but got %s", expect, actual)
			}
			return nil
		}
	}
}

func TestAccEmailSMTPProviderActivation(t *testing.T) {
	frame := test_utils.NewInstanceTestFrame(t, "zitadel_email_provider_smtp")

	initialConfig := fmt.Sprintf(`
%s
resource "zitadel_email_provider_smtp" "default" {
  sender_address = "zitadel@%s"
  sender_name    = "ZITADEL"
  host           = "localhost:25"
  set_active     = false
}
`, frame.ProviderSnippet, frame.InstanceDomain)

	activatedConfig := fmt.Sprintf(`
%s
resource "zitadel_email_provider_smtp" "default" {
  sender_address = "zitadel@%s"
  sender_name    = "ZITADEL"
  host           = "localhost:25"
  set_active     = true
}
`, frame.ProviderSnippet, frame.InstanceDomain)

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

func TestAccEmailSMTPProviderActivationDrift(t *testing.T) {
	frame := test_utils.NewInstanceTestFrame(t, "zitadel_email_provider_smtp")

	activatedConfig := fmt.Sprintf(`
%s
resource "zitadel_email_provider_smtp" "default" {
  sender_address = "zitadel@%s"
  sender_name    = "ZITADEL"
  host           = "localhost:25"
  set_active     = true
}
`, frame.ProviderSnippet, frame.InstanceDomain)

	var providerID string
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: frame.V6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: activatedConfig,
				Check: resource.ComposeTestCheckFunc(
					rememberID(frame, &providerID),
					checkRemoteActive(frame, true),
				),
			},
			{
				PreConfig: func() {
					if _, err := frame.DeactivateEmailProvider(frame, &admin.DeactivateEmailProviderRequest{Id: providerID}); err != nil {
						t.Fatalf("deactivating email provider out of band failed: %v", err)
					}
				},
				Config: activatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(frame.TerraformName, "set_active", "true"),
					checkRemoteActive(frame, true),
				),
			},
		},
	})
}

func TestAccEmailSMTPProviderActivationUnmanaged(t *testing.T) {
	frame := test_utils.NewInstanceTestFrame(t, "zitadel_email_provider_smtp")

	unmanagedConfig := fmt.Sprintf(`
%s
resource "zitadel_email_provider_smtp" "default" {
  sender_address = "zitadel@%s"
  sender_name    = "ZITADEL"
  host           = "localhost:25"
}
`, frame.ProviderSnippet, frame.InstanceDomain)

	var providerID string
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: frame.V6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: unmanagedConfig,
				Check:  rememberID(frame, &providerID),
			},
			{
				PreConfig: func() {
					if _, err := frame.ActivateEmailProvider(frame, &admin.ActivateEmailProviderRequest{Id: providerID}); err != nil {
						t.Fatalf("activating email provider out of band failed: %v", err)
					}
				},
				Config: unmanagedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(frame.TerraformName, "set_active", "true"),
					checkRemoteActive(frame, true),
				),
			},
		},
	})
}

func TestAccEmailSMTPProviderDeactivation(t *testing.T) {
	frame := test_utils.NewInstanceTestFrame(t, "zitadel_email_provider_smtp")

	activatedConfig := fmt.Sprintf(`
%s
resource "zitadel_email_provider_smtp" "default" {
  sender_address = "zitadel@%s"
  sender_name    = "ZITADEL"
  host           = "localhost:25"
  set_active     = true
}
`, frame.ProviderSnippet, frame.InstanceDomain)

	deactivatedConfig := fmt.Sprintf(`
%s
resource "zitadel_email_provider_smtp" "default" {
  sender_address = "zitadel@%s"
  sender_name    = "ZITADEL"
  host           = "localhost:25"
  set_active     = false
}
`, frame.ProviderSnippet, frame.InstanceDomain)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: frame.V6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: activatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(frame.TerraformName, "set_active", "true"),
					checkRemoteActive(frame, true),
				),
			},
			{
				Config: deactivatedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(frame.TerraformName, "set_active", "false"),
					checkRemoteActive(frame, false),
				),
			},
		},
	})
}

func rememberID(frame *test_utils.InstanceTestFrame, id *string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		*id = frame.State(state).ID
		return nil
	}
}

func checkRemoteActive(frame *test_utils.InstanceTestFrame, expect bool) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		resp, err := frame.GetEmailProviderById(frame, &admin.GetEmailProviderByIdRequest{Id: frame.State(state).ID})
		if err != nil {
			return fmt.Errorf("getting email provider failed: %w", err)
		}
		actual := resp.GetConfig().GetState() == settings.EmailProviderState_EMAIL_PROVIDER_ACTIVE
		if actual != expect {
			return fmt.Errorf("expected active %t, but got %t", expect, actual)
		}
		return nil
	}
}
