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
	client, err := helper.GetWebKeyClient(frame, frame.ClientInfo)
	if err != nil {
		t.Fatalf("failed to get client: %v", err)
	}
	// The active key can neither be deleted nor deactivated, so the test hands
	// activation over to a key it does not manage before its keys are destroyed.
	unmanaged, err := client.CreateWebKey(frame, &webkey.CreateWebKeyRequest{})
	if err != nil {
		t.Fatalf("failed to create unmanaged web key: %v", err)
	}
	unmanagedKeyID := unmanaged.GetId()

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

resource "terraform_data" "active_id" {
  input = zitadel_active_webkey.default.id
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

resource "terraform_data" "active_id" {
  input = zitadel_active_webkey.default.id
}
`, frame.ProviderSnippet, frame.AsOrgDefaultDependency)

	configUnmanaged := fmt.Sprintf(`
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
  key_id = "%s"
}

resource "terraform_data" "active_id" {
  input = zitadel_active_webkey.default.id
}
`, frame.ProviderSnippet, frame.AsOrgDefaultDependency, unmanagedKeyID)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: frame.V6ProviderFactories(),
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				Config: configInitial,
				Check: resource.ComposeTestCheckFunc(
					checkRemoteProperty(frame, "key_v1")(""),
					resource.TestCheckResourceAttrPair("terraform_data.active_id", "output", frame.TerraformName, "id"),
				),
			},
			{
				Config: configRotated,
				Check: resource.ComposeTestCheckFunc(
					checkRemoteProperty(frame, "key_v2")(""),
					resource.TestCheckResourceAttrPair("terraform_data.active_id", "output", frame.TerraformName, "id"),
				),
			},
			{
				Config: configUnmanaged,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(frame.TerraformName, "id", frame.OrgID+":"+unmanagedKeyID),
					resource.TestCheckResourceAttrPair("terraform_data.active_id", "output", frame.TerraformName, "id"),
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
