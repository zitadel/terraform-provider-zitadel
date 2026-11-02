package default_security_settings_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	settingsv2 "github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/settings/v2"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper/test_utils"
)

func TestAccDefaultSecuritySettings(t *testing.T) {
	frame := test_utils.NewInstanceTestFrame(t, "zitadel_default_security_settings")
	resourceExample := `
resource "zitadel_default_security_settings" "default" {
  enable_impersonation = true
}
`
	updatedExample := `
resource "zitadel_default_security_settings" "default" {
  enable_impersonation = false
}
`
	test_utils.RunLifecyleTest(
		t,
		frame.BaseTestFrame,
		nil,
		func(property bool, secret string) string {
			if property {
				return resourceExample
			}
			return updatedExample
		},
		true, false,
		"", "", "",
		false,
		checkRemoteProperty(frame),
		regexp.MustCompile(`^default_security_settings$`),
		test_utils.CheckNothing,
		test_utils.ImportNothing,
	)
}

func TestAccDefaultSecuritySettingsCreateZeroValues(t *testing.T) {
	frame := test_utils.NewInstanceTestFrame(t, "zitadel_default_security_settings")

	resourceConfig := fmt.Sprintf(`
%s
resource "zitadel_default_security_settings" "default" {
  enable_impersonation = false
}
`, frame.ProviderSnippet)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: frame.V6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					client, err := helper.GetSecuritySettingsClient(context.Background(), frame.ClientInfo)
					if err != nil {
						t.Fatalf("failed to get client: %v", err)
					}
					if _, err := client.SetSecuritySettings(context.Background(), &settingsv2.SetSecuritySettingsRequest{
						EnableImpersonation: true,
					}); err != nil && helper.IgnorePreconditionError(err) != nil {
						t.Fatalf("setting remote security settings failed: %v", err)
					}
				},
				Config: resourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(frame.TerraformName, "enable_impersonation", "false"),
					test_utils.CheckAMinute(checkRemoteProperty(frame)(false)),
				),
			},
		},
	})
}

// TestAccDefaultSecuritySettingsDynamicClientRegistration checks that the dynamic
// client registration settings are managed, and kept when they are not configured:
// SetSecuritySettings replaces the whole setting, so updating another attribute
// must not reset them.
func TestAccDefaultSecuritySettingsDynamicClientRegistration(t *testing.T) {
	frame := test_utils.NewInstanceTestFrame(t, "zitadel_default_security_settings")

	config := func(body string) string {
		return fmt.Sprintf(`
%s
resource "zitadel_default_security_settings" "default" {
%s
}
`, frame.ProviderSnippet, body)
	}

	allowUnauthenticatedOffConfig := config(`  enable_impersonation                              = false
  embedded_iframe_enabled                           = true
  dynamic_client_registration_allow_unauthenticated = false`)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: frame.V6ProviderFactories(),
		Steps: []resource.TestStep{
			{ // Creating the resource without the attributes keeps the remote values
				PreConfig: func() {
					client, err := helper.GetSecuritySettingsClient(context.Background(), frame.ClientInfo)
					if err != nil {
						t.Fatalf("failed to get client: %v", err)
					}
					if _, err := client.SetSecuritySettings(context.Background(), &settingsv2.SetSecuritySettingsRequest{
						DynamicClientRegistration: &settingsv2.DynamicClientRegistrationSettings{
							Enabled:              true,
							AllowUnauthenticated: true,
						},
					}); err != nil && helper.IgnorePreconditionError(err) != nil {
						t.Fatalf("setting remote security settings failed: %v", err)
					}
				},
				Config: config(`  enable_impersonation = true`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(frame.TerraformName, "dynamic_client_registration_enabled", "true"),
					resource.TestCheckResourceAttr(frame.TerraformName, "dynamic_client_registration_allow_unauthenticated", "true"),
					test_utils.CheckAMinute(checkRemoteDynamicClientRegistration(frame, true, true)),
				),
			},
			{ // Updating other attributes keeps the remote values
				Config: config(`  enable_impersonation    = false
  embedded_iframe_enabled = true`),
				Check: resource.ComposeTestCheckFunc(
					test_utils.CheckAMinute(checkRemoteProperty(frame)(false)),
					test_utils.CheckAMinute(checkRemoteDynamicClientRegistration(frame, true, true)),
				),
			},
			{ // The attributes are set independently
				Config: allowUnauthenticatedOffConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(frame.TerraformName, "dynamic_client_registration_enabled", "true"),
					test_utils.CheckAMinute(checkRemoteDynamicClientRegistration(frame, true, false)),
				),
			},
			{ // Importing reads the attributes
				Config:            allowUnauthenticatedOffConfig,
				ResourceName:      frame.TerraformName,
				ImportState:       true,
				ImportStateId:     "default",
				ImportStateVerify: true,
			},
			{ // Disabling dynamic client registration
				Config: config(`  enable_impersonation                              = false
  dynamic_client_registration_enabled               = false
  dynamic_client_registration_allow_unauthenticated = false`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(frame.TerraformName, "dynamic_client_registration_enabled", "false"),
					test_utils.CheckAMinute(checkRemoteDynamicClientRegistration(frame, false, false)),
				),
			},
		},
	})
}

func checkRemoteDynamicClientRegistration(frame *test_utils.InstanceTestFrame, enabled, allowUnauthenticated bool) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		client, err := helper.GetSecuritySettingsClient(context.Background(), frame.ClientInfo)
		if err != nil {
			return fmt.Errorf("failed to get client: %w", err)
		}
		resp, err := client.GetSecuritySettings(context.Background(), &settingsv2.GetSecuritySettingsRequest{})
		if err != nil {
			return fmt.Errorf("getting security settings failed: %w", err)
		}
		dcr := resp.GetSettings().GetDynamicClientRegistration()
		if dcr.GetEnabled() != enabled || dcr.GetAllowUnauthenticated() != allowUnauthenticated {
			return fmt.Errorf("expected dynamic client registration enabled=%t allow_unauthenticated=%t, but got enabled=%t allow_unauthenticated=%t",
				enabled, allowUnauthenticated, dcr.GetEnabled(), dcr.GetAllowUnauthenticated())
		}
		return nil
	}
}

func checkRemoteProperty(frame *test_utils.InstanceTestFrame) func(bool) resource.TestCheckFunc {
	return func(expect bool) resource.TestCheckFunc {
		return func(state *terraform.State) error {
			client, err := helper.GetSecuritySettingsClient(context.Background(), frame.ClientInfo)
			if err != nil {
				return fmt.Errorf("failed to get client: %w", err)
			}
			resp, err := client.GetSecuritySettings(context.Background(), &settingsv2.GetSecuritySettingsRequest{})
			if err != nil {
				return fmt.Errorf("getting security settings failed: %w", err)
			}
			actual := resp.GetSettings().GetEnableImpersonation()
			if actual != expect {
				return fmt.Errorf("expected %t, but got %t", expect, actual)
			}
			return nil
		}
	}
}
