package webkey_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/webkey/v2"
	"google.golang.org/grpc/metadata"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper/test_utils"
)

func TestAccWebKeyRSA(t *testing.T) {
	frame := test_utils.NewOrgTestFrame(t, "zitadel_webkey")
	config := fmt.Sprintf(`
%s
%s
resource "zitadel_webkey" "default" {
  org_id = data.zitadel_org.default.id
  rsa {
    bits = "RSA_BITS_2048"
  }
}`, frame.ProviderSnippet, frame.AsOrgDefaultDependency)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: frame.V6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					test_utils.CheckStateHasIDSet(frame.BaseTestFrame, test_utils.ZitadelGeneratedIdOnlyRegex),
					checkRemoteProperty(frame, "RSA_BITS_2048", "RSA")(""),
				),
			},
			{
				ResourceName:            frame.TerraformName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"rsa", "ecdsa", "ed25519"},
				ImportStateIdFunc: test_utils.ChainImportStateIdFuncs(
					test_utils.ImportResourceId(frame.BaseTestFrame),
					test_utils.ImportOrgId(frame),
				),
			},
		},
	})
}

func TestAccWebKeyRSA3072(t *testing.T) {
	frame := test_utils.NewOrgTestFrame(t, "zitadel_webkey")
	config := fmt.Sprintf(`
%s
%s
resource "zitadel_webkey" "default" {
  org_id = data.zitadel_org.default.id
  rsa {
    bits = "RSA_BITS_3072"
    hasher = "RSA_HASHER_SHA384"
  }
}`, frame.ProviderSnippet, frame.AsOrgDefaultDependency)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: frame.V6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					test_utils.CheckStateHasIDSet(frame.BaseTestFrame, test_utils.ZitadelGeneratedIdOnlyRegex),
					checkRemoteProperty(frame, "RSA_BITS_3072", "RSA")(""),
				),
			},
		},
	})
}

func TestAccWebKeyRSA4096(t *testing.T) {
	frame := test_utils.NewOrgTestFrame(t, "zitadel_webkey")
	config := fmt.Sprintf(`
%s
%s
resource "zitadel_webkey" "default" {
  org_id = data.zitadel_org.default.id
  rsa {
    bits = "RSA_BITS_4096"
    hasher = "RSA_HASHER_SHA512"
  }
}`, frame.ProviderSnippet, frame.AsOrgDefaultDependency)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: frame.V6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					test_utils.CheckStateHasIDSet(frame.BaseTestFrame, test_utils.ZitadelGeneratedIdOnlyRegex),
					checkRemoteProperty(frame, "RSA_BITS_4096", "RSA")(""),
				),
			},
		},
	})
}

func TestAccWebKeyECDSA(t *testing.T) {
	frame := test_utils.NewOrgTestFrame(t, "zitadel_webkey")
	config := fmt.Sprintf(`
%s
%s
resource "zitadel_webkey" "default" {
  org_id = data.zitadel_org.default.id
  ecdsa {
    curve = "ECDSA_CURVE_P256"
  }
}`, frame.ProviderSnippet, frame.AsOrgDefaultDependency)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: frame.V6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					test_utils.CheckStateHasIDSet(frame.BaseTestFrame, test_utils.ZitadelGeneratedIdOnlyRegex),
					checkRemoteProperty(frame, "ECDSA_CURVE_P256", "ECDSA")(""),
				),
			},
		},
	})
}

func TestAccWebKeyECDSAP384(t *testing.T) {
	frame := test_utils.NewOrgTestFrame(t, "zitadel_webkey")
	config := fmt.Sprintf(`
%s
%s
resource "zitadel_webkey" "default" {
  org_id = data.zitadel_org.default.id
  ecdsa {
    curve = "ECDSA_CURVE_P384"
  }
}`, frame.ProviderSnippet, frame.AsOrgDefaultDependency)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: frame.V6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					test_utils.CheckStateHasIDSet(frame.BaseTestFrame, test_utils.ZitadelGeneratedIdOnlyRegex),
					checkRemoteProperty(frame, "ECDSA_CURVE_P384", "ECDSA")(""),
				),
			},
		},
	})
}

func TestAccWebKeyECDSAP512(t *testing.T) {
	frame := test_utils.NewOrgTestFrame(t, "zitadel_webkey")
	config := fmt.Sprintf(`
%s
%s
resource "zitadel_webkey" "default" {
  org_id = data.zitadel_org.default.id
  ecdsa {
    curve = "ECDSA_CURVE_P512"
  }
}`, frame.ProviderSnippet, frame.AsOrgDefaultDependency)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: frame.V6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					test_utils.CheckStateHasIDSet(frame.BaseTestFrame, test_utils.ZitadelGeneratedIdOnlyRegex),
					checkRemoteProperty(frame, "ECDSA_CURVE_P512", "ECDSA")(""),
				),
			},
		},
	})
}

func TestAccWebKeyED25519(t *testing.T) {
	frame := test_utils.NewOrgTestFrame(t, "zitadel_webkey")
	config := fmt.Sprintf(`
%s
%s
resource "zitadel_webkey" "default" {
  org_id = data.zitadel_org.default.id
  ed25519 {}
}`, frame.ProviderSnippet, frame.AsOrgDefaultDependency)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: frame.V6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					test_utils.CheckStateHasIDSet(frame.BaseTestFrame, test_utils.ZitadelGeneratedIdOnlyRegex),
					checkRemoteProperty(frame, "", "ED25519")(""),
				),
			},
		},
	})
}

func checkRemoteProperty(frame *test_utils.OrgTestFrame, expectedBitsOrCurve, expectedKeyType string) func(string) resource.TestCheckFunc {
	return func(_ string) resource.TestCheckFunc {
		return func(state *terraform.State) error {
			rs, ok := state.RootModule().Resources[frame.TerraformName]
			if !ok {
				return fmt.Errorf("not found: %s", frame.TerraformName)
			}

			client, err := helper.GetWebKeyClient(context.Background(), frame.ClientInfo)
			if err != nil {
				return err
			}

			ctx := metadata.AppendToOutgoingContext(context.Background(), "x-zitadel-orgid", frame.OrgID)
			resp, err := client.ListWebKeys(ctx, &webkey.ListWebKeysRequest{})
			if err != nil {
				return err
			}

			for _, key := range resp.GetWebKeys() {
				if key.GetId() == rs.Primary.ID {
					if expectedKeyType == "RSA" {
						if rsa := key.GetRsa(); rsa != nil {
							if expectedBitsOrCurve != "" && rsa.GetBits().String() != expectedBitsOrCurve {
								return fmt.Errorf("expected rsa bits %s, but got %s", expectedBitsOrCurve, rsa.GetBits().String())
							}
							return nil
						}
						return fmt.Errorf("expected rsa key, but got something else")
					}

					if expectedKeyType == "ECDSA" {
						if ecdsa := key.GetEcdsa(); ecdsa != nil {
							if expectedBitsOrCurve != "" && ecdsa.GetCurve().String() != expectedBitsOrCurve {
								return fmt.Errorf("expected ecdsa curve %s, but got %s", expectedBitsOrCurve, ecdsa.GetCurve().String())
							}
							return nil
						}
						return fmt.Errorf("expected ecdsa key, but got something else")
					}

					if expectedKeyType == "ED25519" {
						if key.GetEd25519() != nil {
							return nil
						}
						return fmt.Errorf("expected ed25519 key, but got something else")
					}

					return nil
				}
			}
			return fmt.Errorf("key with ID %s not found on remote", rs.Primary.ID)
		}
	}
}
