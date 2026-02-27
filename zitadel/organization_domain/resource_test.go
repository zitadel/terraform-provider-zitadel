package organization_domain_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	org "github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/org/v2"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper/test_utils"
)

func TestAccOrganizationDomain(t *testing.T) {
	frame := test_utils.NewInstanceTestFrame(t, "zitadel_organization_domain")
	domainName := frame.UniqueResourcesID + ".example.com"

	resourceConfig := fmt.Sprintf(`
%s
resource "zitadel_organization" "default" {
  name = "%s"
}

resource "zitadel_organization_domain" "default" {
  organization_id = zitadel_organization.default.id
  domain          = "%s"
  validation_type = "DOMAIN_VALIDATION_TYPE_HTTP"
}
`, frame.ProviderSnippet, frame.UniqueResourcesID, domainName)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: frame.V6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: resourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(frame.TerraformName, "domain", domainName),
					resource.TestCheckResourceAttr(frame.TerraformName, "is_verified", "false"),
					resource.TestCheckResourceAttr(frame.TerraformName, "is_primary", "false"),
					resource.TestCheckResourceAttrSet(frame.TerraformName, "validation_token"),
					checkDomainExists(frame, domainName),
				),
			},
		},
	})
}

func checkDomainExists(frame *test_utils.InstanceTestFrame, expectedDomain string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[frame.TerraformName]
		if !ok {
			return fmt.Errorf("not found: %s", frame.TerraformName)
		}

		orgID := rs.Primary.Attributes["organization_id"]
		client, err := helper.GetOrgClient(context.Background(), frame.ClientInfo)
		if err != nil {
			return fmt.Errorf("failed to get client: %w", err)
		}

		resp, err := client.ListOrganizationDomains(context.Background(), &org.ListOrganizationDomainsRequest{
			OrganizationId: orgID,
			Filters: []*org.DomainSearchFilter{
				{
					Filter: &org.DomainSearchFilter_DomainFilter{
						DomainFilter: &org.OrganizationDomainQuery{
							Domain: expectedDomain,
						},
					},
				},
			},
		})
		if err != nil {
			return fmt.Errorf("failed to list domains: %w", err)
		}

		if len(resp.Domains) == 0 {
			return fmt.Errorf("domain %q not found in organization %q", expectedDomain, orgID)
		}

		if resp.Domains[0].Domain != expectedDomain {
			return fmt.Errorf("expected domain %q, but got %q", expectedDomain, resp.Domains[0].Domain)
		}

		return nil
	}
}
