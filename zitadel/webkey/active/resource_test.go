package active_webkey_test

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	webkeys "github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/webkey/v2beta"
	"google.golang.org/grpc/metadata"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper/test_utils"
)

func TestAccActiveWebKey(t *testing.T) {
	frame := test_utils.NewOrgTestFrame(t, "zitadel_active_webkey")
	configInitial := generateResourceHCL(frame, "zitadel_webkey.key_v1.id")
	configRotated := generateResourceHCL(frame, "zitadel_webkey.key_v2.id")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: frame.V6ProviderFactories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: configInitial,
				Check: resource.ComposeTestCheckFunc(
					checkRemoteProperty(frame, "zitadel_webkey.key_v1.id"),
				),
			},
			{
				Config: configRotated,
				Check: resource.ComposeTestCheckFunc(
					checkRemoteProperty(frame, "zitadel_webkey.key_v2.id"),
				),
			},
		},
	})
}

func generateResourceHCL(frame *test_utils.OrgTestFrame, activeKeyRef string) string {
	return fmt.Sprintf(`
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
  key_id = %s
}
`, frame.ProviderSnippet, frame.AsOrgDefaultDependency, activeKeyRef)
}

func checkRemoteProperty(frame *test_utils.OrgTestFrame, expectedKeyIdRef string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		resourceAddress := strings.TrimSuffix(expectedKeyIdRef, ".id")
		rs, ok := state.RootModule().Resources[resourceAddress]
		if !ok {
			return fmt.Errorf("not found in state: %s", resourceAddress)
		}
		expectedKeyId := rs.Primary.ID

		client, err := helper.GetWebKeyClient(context.Background(), frame.ClientInfo)
		if err != nil {
			return err
		}
		ctx := metadata.AppendToOutgoingContext(context.Background(), "x-zitadel-orgid", frame.OrgID)
		resp, err := client.ListWebKeys(ctx, &webkeys.ListWebKeysRequest{})
		if err != nil {
			return err
		}
		for _, key := range resp.GetWebKeys() {
			if key.GetState() == webkeys.State_STATE_ACTIVE {
				if key.GetId() == expectedKeyId {
					return nil
				}
				return fmt.Errorf("expected active key id %s, but got %s", expectedKeyId, key.GetId())
			}
		}
		return fmt.Errorf("no active key found in remote")
	}
}

func checkDestroy(s *terraform.State) error {
	frame := test_utils.NewOrgTestFrame(nil, "zitadel_active_webkey")

	return resource.Retry(15*time.Second, func() *resource.RetryError {
		client, err := helper.GetWebKeyClient(context.Background(), frame.ClientInfo)
		if err != nil {
			return resource.NonRetryableError(err)
		}
		ctx := metadata.AppendToOutgoingContext(context.Background(), "x-zitadel-orgid", frame.OrgID)
		resp, err := client.ListWebKeys(ctx, &webkeys.ListWebKeysRequest{})
		if err != nil {
			return resource.RetryableError(fmt.Errorf("API error while listing keys: %w", err))
		}

		var activeKeyId string
		var oldestKey *webkeys.WebKey

		for _, key := range resp.GetWebKeys() {
			if key.GetState() == webkeys.State_STATE_ACTIVE {
				activeKeyId = key.GetId()
			}

			if oldestKey == nil {
				oldestKey = key
				continue
			}
			if key.GetCreationDate().AsTime().Before(oldestKey.GetCreationDate().AsTime()) {
				oldestKey = key
				continue
			}
			if key.GetCreationDate().AsTime().Equal(oldestKey.GetCreationDate().AsTime()) && key.GetId() < oldestKey.GetId() {
				oldestKey = key
			}
		}

		if oldestKey != nil && activeKeyId != oldestKey.GetId() {
			return resource.RetryableError(fmt.Errorf("destroy did not revert to the oldest key yet: active key is %s, but should be %s", activeKeyId, oldestKey.GetId()))
		}

		return nil
	})
}
