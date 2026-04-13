package sms_provider_http_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper/test_utils"
)

// TestAccSMSHttpProviderSigningKeyRotation reproduces a bug where the update
// function discards the new signing key returned by the API after key rotation.
// After setting expiration_signing_key, the signing_key in state should change.
func TestAccSMSHttpProviderSigningKeyRotation(t *testing.T) {
	frame := test_utils.NewInstanceTestFrame(t, "zitadel_sms_provider_http")

	// Step 1: Create the resource (captures initial signing_key)
	initialConfig := fmt.Sprintf(`
%s
resource "zitadel_sms_provider_http" "default" {
  endpoint = "https://relay.example.com/sms"
}
`, frame.ProviderSnippet)

	// Step 2: Trigger key rotation by changing endpoint (to force an update)
	// and adding expiration_signing_key = "0s" (rotate now, no grace period)
	rotationConfig := fmt.Sprintf(`
%s
resource "zitadel_sms_provider_http" "default" {
  endpoint               = "https://relay.example.com/sms-updated"
  expiration_signing_key = "0s"
}
`, frame.ProviderSnippet)

	var initialKey string

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: frame.V6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: initialConfig,
				Check: resource.ComposeTestCheckFunc(
					// Capture the initial signing key
					func(state *terraform.State) error {
						rs, ok := state.RootModule().Resources[frame.TerraformName]
						if !ok {
							return fmt.Errorf("resource not found: %s", frame.TerraformName)
						}
						initialKey = rs.Primary.Attributes["signing_key"]
						if initialKey == "" {
							return fmt.Errorf("signing_key is empty after create")
						}
						t.Logf("Initial signing_key: %s", initialKey[:10]+"...")
						return nil
					},
				),
			},
			{
				Config: rotationConfig,
				Check: resource.ComposeTestCheckFunc(
					// After rotation, signing_key should be different
					func(state *terraform.State) error {
						rs, ok := state.RootModule().Resources[frame.TerraformName]
						if !ok {
							return fmt.Errorf("resource not found: %s", frame.TerraformName)
						}
						newKey := rs.Primary.Attributes["signing_key"]
						if newKey == "" {
							return fmt.Errorf("signing_key is empty after rotation")
						}
						if newKey == initialKey {
							return fmt.Errorf("signing_key did not change after rotation: still %s", newKey[:10]+"...")
						}
						t.Logf("New signing_key: %s", newKey[:10]+"...")
						return nil
					},
				),
			},
		},
	})
}
