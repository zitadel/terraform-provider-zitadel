package webkey_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper/test_utils"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/webkey"
)

func TestAccWebKeyDatasource(t *testing.T) {
	frame := test_utils.NewOrgTestFrame(t, "zitadel_webkey_datasource")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: frame.V6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
%s
%s

resource "zitadel_webkey" "testkey" {
  org_id = data.zitadel_org.default.id
  rsa {}
}

data "zitadel_webkey" "default" {
  org_id    = data.zitadel_org.default.id
  webkey_id = zitadel_webkey.testkey.id
}
`, frame.ProviderSnippet, frame.AsOrgDefaultDependency),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						"data.zitadel_webkey.default", "id",
						"zitadel_webkey.testkey", "id",
					),
					resource.TestCheckResourceAttr("data.zitadel_webkey.default", webkey.KeyTypeVar, "RSA"),
					resource.TestCheckResourceAttrSet("data.zitadel_webkey.default", webkey.StateVar),
				),
			},
		},
	})
}
