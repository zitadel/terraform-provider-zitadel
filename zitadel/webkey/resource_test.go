package webkey_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/webkey/v2beta"
	"google.golang.org/grpc/metadata"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper/test_utils"
)

func TestAccWebKey(t *testing.T) {
	frame := test_utils.NewOrgTestFrame(t, "zitadel_webkey")
	config := generateResourceHCL(frame, "RSA_BITS_2048")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: frame.V6ProviderFactories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					test_utils.CheckStateHasIDSet(frame.BaseTestFrame, test_utils.ZitadelGeneratedIdOnlyRegex),
					checkRemoteProperty(frame, "RSA_BITS_2048"),
				),
			},
		},
	})
}

func generateResourceHCL(frame *test_utils.OrgTestFrame, bits string) string {
	return fmt.Sprintf(`
%s
%s
resource "zitadel_webkey" "default" {
  org_id = data.zitadel_org.default.id
  rsa {
    bits = "%s"
  }
}`, frame.ProviderSnippet, frame.AsOrgDefaultDependency, bits)
}

func checkRemoteProperty(frame *test_utils.OrgTestFrame, expectedBits string) resource.TestCheckFunc {
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
				if rsa := key.GetRsa(); rsa != nil {
					if rsa.GetBits().String() == expectedBits {
						return nil
					}
					return fmt.Errorf("expected rsa bits %s, but got %s", expectedBits, rsa.GetBits().String())
				}
				return fmt.Errorf("expected rsa key, but got something else")
			}
		}
		return fmt.Errorf("key with ID %s not found on remote", rs.Primary.ID)
	}
}

func checkDestroy(s *terraform.State) error {
	frame := test_utils.NewOrgTestFrame(nil, "zitadel_webkey")
	client, err := helper.GetWebKeyClient(context.Background(), frame.ClientInfo)
	if err != nil {
		return nil
	}
	ctx := metadata.AppendToOutgoingContext(context.Background(), "x-zitadel-orgid", frame.OrgID)
	resp, err := client.ListWebKeys(ctx, &webkey.ListWebKeysRequest{})
	if err != nil {
		return nil
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "zitadel_webkey" {
			continue
		}
		for _, key := range resp.GetWebKeys() {
			if key.GetId() == rs.Primary.ID {
				return fmt.Errorf("webkey with id %s still exists", rs.Primary.ID)
			}
		}
	}
	return nil
}
