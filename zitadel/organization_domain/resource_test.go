package organization_domain_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	org "github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/org/v2"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper/test_utils"
)

// domainPolicy returns an org level domain policy that either requires org
// domains to be validated or lets ZITADEL verify them on add. The acceptance
// instances inherit the instance default, so every test pins the policy it
// needs instead of relying on it.
func domainPolicy(validateOrgDomains bool) string {
	return fmt.Sprintf(`
resource "zitadel_domain_policy" "default" {
  org_id                                      = zitadel_organization.default.id
  user_login_must_be_domain                   = false
  validate_org_domains                        = %t
  smtp_sender_address_matches_instance_domain = false
}`, validateOrgDomains)
}

// TestAccOrganizationDomain covers the default domain policy, where
// validate_org_domains is disabled and ZITADEL verifies the domain as part of
// adding it. There is nothing to validate, so validation_type can be omitted
// and the resource must still land in state (issue #427).
func TestAccOrganizationDomain(t *testing.T) {
	frame := test_utils.NewInstanceTestFrame(t, "zitadel_organization_domain")
	domainName := frame.UniqueResourcesID + ".example.com"

	resourceConfig := fmt.Sprintf(`
%s
resource "zitadel_organization" "default" {
  name = "%s"
}
%s

resource "zitadel_organization_domain" "default" {
  organization_id = zitadel_organization.default.id
  domain          = "%s"
  depends_on      = [zitadel_domain_policy.default]
}
`, frame.ProviderSnippet, frame.UniqueResourcesID, domainPolicy(false), domainName)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: frame.V6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: resourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(frame.TerraformName, "domain", domainName),
					resource.TestCheckResourceAttr(frame.TerraformName, "is_verified", "true"),
					resource.TestCheckResourceAttr(frame.TerraformName, "is_primary", "false"),
					checkDomainExists(frame, domainName),
				),
			},
		},
	})
}

// TestAccOrganizationDomainAutoVerifiedWithValidationType is the migration
// path from the deprecated zitadel_domain resource: validation_type is set,
// but the policy auto-verifies the domain anyway. Generating a validation
// challenge would fail with ORG-HGw21, so it has to be skipped rather than
// aborting the create.
func TestAccOrganizationDomainAutoVerifiedWithValidationType(t *testing.T) {
	frame := test_utils.NewInstanceTestFrame(t, "zitadel_organization_domain")
	domainName := frame.UniqueResourcesID + ".example.com"

	resourceConfig := fmt.Sprintf(`
%s
resource "zitadel_organization" "default" {
  name = "%s"
}
%s

resource "zitadel_organization_domain" "default" {
  organization_id = zitadel_organization.default.id
  domain          = "%s"
  validation_type = "DOMAIN_VALIDATION_TYPE_HTTP"
  depends_on      = [zitadel_domain_policy.default]
}
`, frame.ProviderSnippet, frame.UniqueResourcesID, domainPolicy(false), domainName)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: frame.V6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: resourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(frame.TerraformName, "domain", domainName),
					resource.TestCheckResourceAttr(frame.TerraformName, "is_verified", "true"),
					resource.TestCheckNoResourceAttr(frame.TerraformName, "validation_token"),
					checkDomainExists(frame, domainName),
				),
			},
		},
	})
}

// TestAccOrganizationDomainWithValidation covers the policy that requires org
// domains to be validated. The domain stays unverified after being added, so a
// validation challenge is generated and exposed.
func TestAccOrganizationDomainWithValidation(t *testing.T) {
	frame := test_utils.NewInstanceTestFrame(t, "zitadel_organization_domain")
	domainName := frame.UniqueResourcesID + ".example.com"

	resourceConfig := fmt.Sprintf(`
%s
resource "zitadel_organization" "default" {
  name = "%s"
}
%s

resource "zitadel_organization_domain" "default" {
  organization_id = zitadel_organization.default.id
  domain          = "%s"
  validation_type = "DOMAIN_VALIDATION_TYPE_HTTP"
  depends_on      = [zitadel_domain_policy.default]
}
`, frame.ProviderSnippet, frame.UniqueResourcesID, domainPolicy(true), domainName)

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

// TestAccOrganizationDomainVerifyWithoutValidationType covers the combination
// that making validation_type optional newly allows to be expressed but that
// can never succeed: under a policy that requires validation there is no
// challenge for ZITADEL to check, so verify is rejected with an actionable
// error before the domain is added.
func TestAccOrganizationDomainVerifyWithoutValidationType(t *testing.T) {
	frame := test_utils.NewInstanceTestFrame(t, "zitadel_organization_domain")
	domainName := frame.UniqueResourcesID + ".example.com"

	resourceConfig := fmt.Sprintf(`
%s
resource "zitadel_organization" "default" {
  name = "%s"
}
%s

resource "zitadel_organization_domain" "default" {
  organization_id = zitadel_organization.default.id
  domain          = "%s"
  verify          = true
  depends_on      = [zitadel_domain_policy.default]
}
`, frame.ProviderSnippet, frame.UniqueResourcesID, domainPolicy(true), domainName)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: frame.V6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config:      resourceConfig,
				ExpectError: regexp.MustCompile(`verify needs validation_type when the domain policy`),
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
