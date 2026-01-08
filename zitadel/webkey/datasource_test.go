package webkey_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper/test_utils"
)

func TestAccWebKeyDatasource(t *testing.T) {
	frame := test_utils.NewOrgTestFrame(t, "zitadel_webkey")
	config := fmt.Sprintf(`
%s
%s
resource "zitadel_webkey" "default" {
  org_id = data.zitadel_org.default.id
  rsa {
    bits = "RSA_BITS_2048"
  }
}

data "zitadel_webkey" "default" {
  org_id    = data.zitadel_org.default.id
  webkey_id = zitadel_webkey.default.id
}
`, frame.ProviderSnippet, frame.AsOrgDefaultDependency)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: frame.V6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.zitadel_webkey.default", "id", "zitadel_webkey.default", "id"),
				),
			},
		},
	})
}

func TestAccWebKeyDatasourceECDSA(t *testing.T) {
	frame := test_utils.NewOrgTestFrame(t, "zitadel_webkey")
	config := fmt.Sprintf(`
%s
%s

resource "zitadel_webkey" "testkey" {
  org_id = data.zitadel_org.default.id
  ecdsa {
    curve = "ECDSA_CURVE_P256"
  }
}

data "zitadel_webkey" "default" {
  org_id    = data.zitadel_org.default.id
  webkey_id = zitadel_webkey.testkey.id
}
`, frame.ProviderSnippet, frame.AsOrgDefaultDependency)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: frame.V6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.zitadel_webkey.default", "id", "zitadel_webkey.testkey", "id"),
				),
			},
		},
	})
}

func TestAccWebKeyDatasourceED25519(t *testing.T) {
	frame := test_utils.NewOrgTestFrame(t, "zitadel_webkey")
	config := fmt.Sprintf(`
%s
%s

resource "zitadel_webkey" "testkey" {
  org_id = data.zitadel_org.default.id
  ed25519 {}
}

data "zitadel_webkey" "default" {
  org_id    = data.zitadel_org.default.id
  webkey_id = zitadel_webkey.testkey.id
}
`, frame.ProviderSnippet, frame.AsOrgDefaultDependency)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: frame.V6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.zitadel_webkey.default", "id", "zitadel_webkey.testkey", "id"),
				),
			},
		},
	})
}
