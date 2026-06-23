package organization_domain_validation_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/ephemeral/ephtest"
)

// TestAccOrganizationDomainValidation verifies the ephemeral resource
// (re)generates the verification token for an organization domain. A fresh org
// and an unverified domain are created first; the ephemeral resource then
// generates a validation challenge for that domain.
func TestAccOrganizationDomainValidation(t *testing.T) {
	// Use the instance-level frame: its domain policy keeps added org domains
	// unverified, so there is an actual validation challenge to (re)generate.
	// The org-level instance auto-verifies domains, leaving nothing to validate.
	frame := ephtest.NewInstanceFrame(t)
	config := fmt.Sprintf(`
%s
resource "zitadel_organization" "default" {
  name = "ephtest-%s"
}
# Force org-domain validation so an added domain stays unverified and there is
# an actual challenge to (re)generate; otherwise the instance auto-verifies it.
resource "zitadel_domain_policy" "default" {
  org_id                                      = zitadel_organization.default.id
  user_login_must_be_domain                   = false
  validate_org_domains                        = true
  smtp_sender_address_matches_instance_domain = false
}
resource "zitadel_organization_domain" "default" {
  organization_id = zitadel_organization.default.id
  domain          = "ephtest-%s.example.com"
  validation_type = "DOMAIN_VALIDATION_TYPE_DNS"
  depends_on      = [zitadel_domain_policy.default]
}
ephemeral "zitadel_organization_domain_validation" "test" {
  org_id          = zitadel_organization.default.id
  domain          = zitadel_organization_domain.default.domain
  validation_type = "DOMAIN_VALIDATION_TYPE_DNS"
}
provider "echo" {
  data = ephemeral.zitadel_organization_domain_validation.test
}
resource "echo" "test" {}
`, frame.ProviderSnippet, frame.UniqueID, frame.UniqueID)

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_10_0),
		},
		ProtoV6ProviderFactories: frame.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("echo.test", "data.token"),
				),
			},
		},
	})
}
