package active_webkey_test

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

func TestAccActiveWebKey(t *testing.T) {
	frame := test_utils.NewOrgTestFrame(t, "zitadel_active_webkey")
	configInitial := fmt.Sprintf(`
%s
%s

resource "zitadel_webkey" "key_v1" {
  org_id = data.zitadel_org.default.id
  rsa {}
}

resource "zitadel_webkey" "key_v2" {
  org_id = data.zitadel_org.default.id
  ecdsa {}
}

resource "zitadel_active_webkey" "default" {
  org_id = data.zitadel_org.default.id
  key_id = zitadel_webkey.key_v1.id
}
`, frame.ProviderSnippet, frame.AsOrgDefaultDependency)

	configRotated := fmt.Sprintf(`
%s
%s

resource "zitadel_webkey" "key_v1" {
  org_id = data.zitadel_org.default.id
  rsa {}
}

resource "zitadel_webkey" "key_v2" {
  org_id = data.zitadel_org.default.id
  ecdsa {}
}

resource "zitadel_active_webkey" "default" {
  org_id = data.zitadel_org.default.id
  key_id = zitadel_webkey.key_v2.id
}
`, frame.ProviderSnippet, frame.AsOrgDefaultDependency)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: frame.V6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: configInitial,
				Check: resource.ComposeTestCheckFunc(
					checkRemoteProperty(frame, "key_v1")(""),
				),
			},
			{
				Config: configRotated,
				Check: resource.ComposeTestCheckFunc(
					checkRemoteProperty(frame, "key_v2")(""),
				),
			},
		},
	})
}

func TestAccActiveWebKeyRotation(t *testing.T) {
	frame := test_utils.NewOrgTestFrame(t, "zitadel_active_webkey")
	configInitial := fmt.Sprintf(`
%s
%s

resource "zitadel_webkey" "key_v1" {
  org_id = data.zitadel_org.default.id
  rsa {
    bits = "RSA_BITS_2048"
  }
}

resource "zitadel_webkey" "key_v2" {
  org_id = data.zitadel_org.default.id
  ecdsa {
    curve = "ECDSA_CURVE_P256"
  }
}

resource "zitadel_active_webkey" "default" {
  org_id = data.zitadel_org.default.id
  key_id = zitadel_webkey.key_v1.id
}
`, frame.ProviderSnippet, frame.AsOrgDefaultDependency)

	configRotated := fmt.Sprintf(`
%s
%s

resource "zitadel_webkey" "key_v1" {
  org_id = data.zitadel_org.default.id
  rsa {
    bits = "RSA_BITS_2048"
  }
}

resource "zitadel_webkey" "key_v2" {
  org_id = data.zitadel_org.default.id
  ecdsa {
    curve = "ECDSA_CURVE_P256"
  }
}

resource "zitadel_active_webkey" "default" {
  org_id = data.zitadel_org.default.id
  key_id = zitadel_webkey.key_v2.id
}
`, frame.ProviderSnippet, frame.AsOrgDefaultDependency)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: frame.V6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: configInitial,
				Check: resource.ComposeTestCheckFunc(
					checkRemoteProperty(frame, "key_v1")(""),
				),
			},
			{
				Config: configRotated,
				Check: resource.ComposeTestCheckFunc(
					checkRemoteProperty(frame, "key_v2")(""),
				),
			},
		},
	})
}

func checkRemoteProperty(frame *test_utils.OrgTestFrame, expectedKeyRef string) func(string) resource.TestCheckFunc {
	return func(_ string) resource.TestCheckFunc {
		return func(state *terraform.State) error {
			rs, ok := state.RootModule().Resources["zitadel_webkey."+expectedKeyRef]
			if !ok {
				return fmt.Errorf("not found in state: zitadel_webkey.%s", expectedKeyRef)
			}
			expectedKeyId := rs.Primary.ID

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
				if key.GetState() == webkey.State_STATE_ACTIVE {
					if key.GetId() == expectedKeyId {
						return nil
					}
					return fmt.Errorf("expected active key id %s, but got %s", expectedKeyId, key.GetId())
				}
			}
			return fmt.Errorf("no active key found in remote")
		}
	}
}
