package organization_metadata_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper/test_utils"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/organization_metadata"
)

func TestAccOrganizationMetadata(t *testing.T) {
	frame := test_utils.NewInstanceTestFrame(t, "zitadel_organization_metadata")
	key := "key_" + frame.UniqueResourcesID

	resourceConfig := fmt.Sprintf(`
%s
resource "zitadel_organization" "default" {
  name = "%s"
}

resource "zitadel_organization_metadata" "default" {
  organization_id = zitadel_organization.default.id
  key             = "%s"
  value           = "example_value"
}
`, frame.ProviderSnippet, frame.UniqueResourcesID, key)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: frame.V6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: resourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(frame.TerraformName, organization_metadata.KeyVar, key),
					resource.TestCheckResourceAttr(frame.TerraformName, organization_metadata.ValueVar, "example_value"),
				),
			},
			{
				// The import ID is <organization_id:key> (issue #452).
				ResourceName: frame.TerraformName,
				ImportState:  true,
				ImportStateIdFunc: test_utils.ChainImportStateIdFuncs(
					test_utils.ImportStateAttribute(frame.BaseTestFrame, organization_metadata.OrganizationIDVar),
					test_utils.ImportStateAttribute(frame.BaseTestFrame, organization_metadata.KeyVar),
				),
				ImportStateVerify: true,
			},
		},
	})
}
