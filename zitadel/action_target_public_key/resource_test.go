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

	resourceConfig := fmt.Sprintf(`
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
-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA0Z3VS5JJcds3xfn/ygWe
FsXpOJFdGMqhBJCnISAAnNPBKSFwETb4FIxgpJMtzBCIR2YEKXE6OryMpO6E8yoI
6sFawwLY1ViELOE7FD7sJVMUQF1WLiMjb7n1feGfToGarnWjKrx8IXjlgVnJ5kQ0
GNOwjKBOmgJiJEhBuTflS0ppODBdKP2oq6iAdf5bMmkv0wMKJnxBKPQsXLcCn2u4
ym9AXkcdH2QviCBWMpGrjVoGLFGqf5E4MiwMuNl7rHIExmBm2mlnmuIPhILRs/jS
tKKLrdazqFCxD2fWXt9a2yzXoE6Hv0sWBnJSRASez2dn6ki3GFbLHeR2dMhT8wbf
cQIDAQAB
-----END PUBLIC KEY-----
EOT
}
`, frame.ProviderSnippet, frame.UniqueResourcesID)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: frame.V6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: resourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(frame.TerraformName, "key_id"),
					test_utils.CheckAMinute(checkRemoteProperty(frame)),
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

func checkRemoteProperty(frame *test_utils.InstanceTestFrame) resource.TestCheckFunc {
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

		return nil
	}
}
