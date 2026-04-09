package action_target_public_key_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper/test_utils"
)

func TestAccActionTargetPublicKeyDatasource(t *testing.T) {
	frame := test_utils.NewInstanceTestFrame(t, "zitadel_action_target_public_key")
	config := fmt.Sprintf(`
%s
resource "zitadel_action_target" "default" {
  name               = "%s"
  endpoint           = "https://example.com/datasource-test"
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

data "zitadel_action_target_public_key" "default" {
  target_id = zitadel_action_target.default.id
  key_id    = zitadel_action_target_public_key.default.id
}
`, frame.ProviderSnippet, frame.UniqueResourcesID)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: frame.V6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.zitadel_action_target_public_key.default", "id", "zitadel_action_target_public_key.default", "id"),
					resource.TestCheckResourceAttrSet("data.zitadel_action_target_public_key.default", "fingerprint"),
					resource.TestCheckResourceAttrSet("data.zitadel_action_target_public_key.default", "public_key"),
				),
			},
		},
	})
}
