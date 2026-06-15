package action_target_public_key_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	actionv2 "github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/action/v2"
	filterv2 "github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/filter/v2"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper/test_utils"
)

func TestAccActionTargetPublicKey(t *testing.T) {
	frame := test_utils.NewInstanceTestFrame(t, "zitadel_action_target_public_key")

	const publicKey = `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA0Z3VS5JJcds3xfn/ygWe
FsXpOJFdGMqhBJCnISAAnNPBKSFwETb4FIxgpJMtzBCIR2YEKXE6OryMpO6E8yoI
6sFawwLY1ViELOE7FD7sJVMUQF1WLiMjb7n1feGfToGarnWjKrx8IXjlgVnJ5kQ0
GNOwjKBOmgJiJEhBuTflS0ppODBdKP2oq6iAdf5bMmkv0wMKJnxBKPQsXLcCn2u4
ym9AXkcdH2QviCBWMpGrjVoGLFGqf5E4MiwMuNl7rHIExmBm2mlnmuIPhILRs/jS
tKKLrdazqFCxD2fWXt9a2yzXoE6Hv0sWBnJSRASez2dn6ki3GFbLHeR2dMhT8wbf
cQIDAQAB
-----END PUBLIC KEY-----`

	configWithoutActive := fmt.Sprintf(`
%s
resource "zitadel_action_target" "default" {
  name               = "%s"
  endpoint           = "https://example.com/test"
  target_type        = "REST_ASYNC"
  timeout            = "10s"
  interrupt_on_error = false
  payload_type       = "PAYLOAD_TYPE_JWE"
}

resource "zitadel_action_target_public_key" "default" {
  target_id  = zitadel_action_target.default.id
  public_key = <<-EOT
%s
EOT
}
`, frame.ProviderSnippet, frame.UniqueResourcesID, publicKey)

	configActive := fmt.Sprintf(`
%s
resource "zitadel_action_target" "default" {
  name               = "%s"
  endpoint           = "https://example.com/test"
  target_type        = "REST_ASYNC"
  timeout            = "10s"
  interrupt_on_error = false
  payload_type       = "PAYLOAD_TYPE_JWE"
}

resource "zitadel_action_target_public_key" "default" {
  target_id  = zitadel_action_target.default.id
  active     = true
  public_key = <<-EOT
%s
EOT
}
`, frame.ProviderSnippet, frame.UniqueResourcesID, publicKey)

	configInactive := fmt.Sprintf(`
%s
resource "zitadel_action_target" "default" {
  name               = "%s"
  endpoint           = "https://example.com/test"
  target_type        = "REST_ASYNC"
  timeout            = "10s"
  interrupt_on_error = false
  payload_type       = "PAYLOAD_TYPE_JWE"
}

resource "zitadel_action_target_public_key" "default" {
  target_id  = zitadel_action_target.default.id
  active     = false
  public_key = <<-EOT
%s
EOT
}
`, frame.ProviderSnippet, frame.UniqueResourcesID, publicKey)

	// Capture the key ID after the initial create so subsequent toggle steps can
	// assert it stays stable — proves activation toggling is an in-place update,
	// not a recreate that would break rotation flows.
	var initialKeyID string
	captureKeyID := func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[frame.TerraformName]
		if !ok {
			return fmt.Errorf("not found: %s", frame.TerraformName)
		}
		initialKeyID = rs.Primary.ID
		return nil
	}
	assertKeyIDStable := func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[frame.TerraformName]
		if !ok {
			return fmt.Errorf("not found: %s", frame.TerraformName)
		}
		if rs.Primary.ID != initialKeyID {
			return fmt.Errorf("key_id changed across toggle: want %q, got %q", initialKeyID, rs.Primary.ID)
		}
		return nil
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: frame.V6ProviderFactories(),
		Steps: []resource.TestStep{
			// Pre-existing behavior: no `active` in config -> key is created and remains inactive.
			{
				Config: configWithoutActive,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(frame.TerraformName, "key_id"),
					resource.TestCheckResourceAttr(frame.TerraformName, "active", "false"),
					captureKeyID,
					test_utils.CheckAMinute(checkRemoteProperty(frame, false)),
				),
			},
			// Toggle to active=true via update (no recreate).
			{
				Config: configActive,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(frame.TerraformName, "active", "true"),
					assertKeyIDStable,
					test_utils.CheckAMinute(checkRemoteProperty(frame, true)),
				),
			},
			// Toggle back to active=false via update.
			{
				Config: configInactive,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(frame.TerraformName, "active", "false"),
					assertKeyIDStable,
					test_utils.CheckAMinute(checkRemoteProperty(frame, false)),
				),
			},
			{
				ResourceName:            frame.TerraformName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"public_key"},
				ImportStateIdFunc: test_utils.ChainImportStateIdFuncs(
					test_utils.ImportResourceId(frame.BaseTestFrame),
					test_utils.ImportStateAttribute(frame.BaseTestFrame, "target_id"),
				),
			},
		},
	})
}

// TestAccActionTargetPublicKeyCreateActive verifies that a resource created with
// active=true is activated on the server during Create (not just after a subsequent
// Update), so the key is usable for payload encryption immediately after apply.
func TestAccActionTargetPublicKeyCreateActive(t *testing.T) {
	frame := test_utils.NewInstanceTestFrame(t, "zitadel_action_target_public_key")

	const publicKey = `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA0Z3VS5JJcds3xfn/ygWe
FsXpOJFdGMqhBJCnISAAnNPBKSFwETb4FIxgpJMtzBCIR2YEKXE6OryMpO6E8yoI
6sFawwLY1ViELOE7FD7sJVMUQF1WLiMjb7n1feGfToGarnWjKrx8IXjlgVnJ5kQ0
GNOwjKBOmgJiJEhBuTflS0ppODBdKP2oq6iAdf5bMmkv0wMKJnxBKPQsXLcCn2u4
ym9AXkcdH2QviCBWMpGrjVoGLFGqf5E4MiwMuNl7rHIExmBm2mlnmuIPhILRs/jS
tKKLrdazqFCxD2fWXt9a2yzXoE6Hv0sWBnJSRASez2dn6ki3GFbLHeR2dMhT8wbf
cQIDAQAB
-----END PUBLIC KEY-----`

	configActiveOnCreate := fmt.Sprintf(`
%s
resource "zitadel_action_target" "default" {
  name               = "%s"
  endpoint           = "https://example.com/test"
  target_type        = "REST_ASYNC"
  timeout            = "10s"
  interrupt_on_error = false
  payload_type       = "PAYLOAD_TYPE_JWE"
}

resource "zitadel_action_target_public_key" "default" {
  target_id  = zitadel_action_target.default.id
  active     = true
  public_key = <<-EOT
%s
EOT
}
`, frame.ProviderSnippet, frame.UniqueResourcesID, publicKey)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: frame.V6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: configActiveOnCreate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(frame.TerraformName, "key_id"),
					resource.TestCheckResourceAttr(frame.TerraformName, "active", "true"),
					test_utils.CheckAMinute(checkRemoteProperty(frame, true)),
				),
			},
		},
	})
}

// TestAccActionTargetPublicKeyNoAccidentalToggle proves the upgrade is non-breaking:
//
//   - A key activated out-of-band (e.g. via the API by a user who upgraded from a
//     provider release that did not manage activation) must NOT be silently
//     deactivated when terraform applies a config that omits the `active` field.
//   - Removing `active` from config must NOT change the remote activation state.
//   - Adding `active = true` to config when the server is already active must be a
//     plan/apply no-op (FailedPrecondition idempotency).
func TestAccActionTargetPublicKeyNoAccidentalToggle(t *testing.T) {
	frame := test_utils.NewInstanceTestFrame(t, "zitadel_action_target_public_key")

	const publicKey = `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA0Z3VS5JJcds3xfn/ygWe
FsXpOJFdGMqhBJCnISAAnNPBKSFwETb4FIxgpJMtzBCIR2YEKXE6OryMpO6E8yoI
6sFawwLY1ViELOE7FD7sJVMUQF1WLiMjb7n1feGfToGarnWjKrx8IXjlgVnJ5kQ0
GNOwjKBOmgJiJEhBuTflS0ppODBdKP2oq6iAdf5bMmkv0wMKJnxBKPQsXLcCn2u4
ym9AXkcdH2QviCBWMpGrjVoGLFGqf5E4MiwMuNl7rHIExmBm2mlnmuIPhILRs/jS
tKKLrdazqFCxD2fWXt9a2yzXoE6Hv0sWBnJSRASez2dn6ki3GFbLHeR2dMhT8wbf
cQIDAQAB
-----END PUBLIC KEY-----`

	target := fmt.Sprintf(`
%s
resource "zitadel_action_target" "default" {
  name               = "%s"
  endpoint           = "https://example.com/test"
  target_type        = "REST_ASYNC"
  timeout            = "10s"
  interrupt_on_error = false
  payload_type       = "PAYLOAD_TYPE_JWE"
}
`, frame.ProviderSnippet, frame.UniqueResourcesID)

	configNoActive := target + fmt.Sprintf(`
resource "zitadel_action_target_public_key" "default" {
  target_id  = zitadel_action_target.default.id
  public_key = <<-EOT
%s
EOT
}
`, publicKey)

	configActiveTrue := target + fmt.Sprintf(`
resource "zitadel_action_target_public_key" "default" {
  target_id  = zitadel_action_target.default.id
  active     = true
  public_key = <<-EOT
%s
EOT
}
`, publicKey)

	var captured struct {
		targetID, keyID string
	}
	captureIDs := func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[frame.TerraformName]
		if !ok {
			return fmt.Errorf("not found: %s", frame.TerraformName)
		}
		captured.targetID = rs.Primary.Attributes["target_id"]
		captured.keyID = rs.Primary.ID
		return nil
	}
	activateExternally := func() {
		if captured.targetID == "" || captured.keyID == "" {
			t.Fatal("captureIDs did not run before activateExternally")
		}
		client, err := helper.GetActionClient(context.Background(), frame.ClientInfo)
		if err != nil {
			t.Fatalf("failed to get client: %v", err)
		}
		if _, err := client.ActivatePublicKey(context.Background(), &actionv2.ActivatePublicKeyRequest{
			TargetId: captured.targetID,
			KeyId:    captured.keyID,
		}); err != nil && helper.IgnorePreconditionError(err) != nil {
			t.Fatalf("external activation failed: %v", err)
		}
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: frame.V6ProviderFactories(),
		Steps: []resource.TestStep{
			// Baseline: create without `active` -> key inactive on server (mirrors the
			// behaviour every existing user has been getting from prior provider releases).
			{
				Config: configNoActive,
				Check: resource.ComposeTestCheckFunc(
					captureIDs,
					resource.TestCheckResourceAttr(frame.TerraformName, "active", "false"),
					test_utils.CheckAMinute(checkRemoteProperty(frame, false)),
				),
			},
			// User (or another tool) activates the key out-of-band via the ZITADEL API.
			// Apply the SAME config: the provider must refresh state, see active=true,
			// and produce no plan diff / no deactivation. This is the critical upgrade
			// safety property.
			{
				PreConfig: activateExternally,
				Config:    configNoActive,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(frame.TerraformName, "active", "true"),
					test_utils.CheckAMinute(checkRemoteProperty(frame, true)),
				),
			},
			// Explicitly verify that, after the out-of-band activation has been
			// reconciled into state, a follow-up plan against the same `active`-less
			// config is a no-op. PlanOnly fails the test if the plan is non-empty,
			// which would catch any regression to "removing `active` from config
			// silently deactivates the key".
			{
				Config:   configNoActive,
				PlanOnly: true,
			},
			// Adding `active = true` to config when the server already matches must be a
			// no-op (idempotent activate; FailedPrecondition is swallowed).
			{
				Config: configActiveTrue,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(frame.TerraformName, "active", "true"),
					test_utils.CheckAMinute(checkRemoteProperty(frame, true)),
				),
			},
			// Removing `active` from config again must NOT deactivate the key. The
			// remote state stays active even though the field is gone from config.
			{
				Config: configNoActive,
				Check: resource.ComposeTestCheckFunc(
					test_utils.CheckAMinute(checkRemoteProperty(frame, true)),
				),
			},
		},
	})
}

func checkRemoteProperty(frame *test_utils.InstanceTestFrame, wantActive bool) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[frame.TerraformName]
		if !ok {
			return fmt.Errorf("not found: %s", frame.TerraformName)
		}

		targetID := rs.Primary.Attributes["target_id"]
		keyID := rs.Primary.ID

		client, err := helper.GetActionClient(context.Background(), frame.ClientInfo)
		if err != nil {
			return fmt.Errorf("failed to get client: %w", err)
		}

		resp, err := client.ListPublicKeys(context.Background(), &actionv2.ListPublicKeysRequest{
			TargetId: targetID,
			Filters: []*actionv2.PublicKeySearchFilter{
				{
					Filter: &actionv2.PublicKeySearchFilter_KeyIdsFilter{
						KeyIdsFilter: &filterv2.InIDsFilter{
							Ids: []string{keyID},
						},
					},
				},
			},
		})
		if err != nil {
			return err
		}

		keys := resp.GetPublicKeys()
		if len(keys) == 0 {
			return fmt.Errorf("public key %s not found on target %s", keyID, targetID)
		}

		if got := keys[0].GetActive(); got != wantActive {
			return fmt.Errorf("public key %s on target %s: want active=%v, got active=%v", keyID, targetID, wantActive, got)
		}

		return nil
	}
}
