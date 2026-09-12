# With the default domain policy (validate_org_domains = false) ZITADEL
# verifies the domain while adding it, so no validation_type is needed.
resource "zitadel_organization_domain" "default" {
  organization_id = zitadel_organization.default.id
  domain          = "example.com"
}

# When the domain policy requires org domains to be validated, pick a
# validation_type to get a challenge to publish.
resource "zitadel_organization_domain" "validated" {
  organization_id = zitadel_organization.default.id
  domain          = "validated.example.com"
  validation_type = "DOMAIN_VALIDATION_TYPE_DNS"
}

output "dns_validation_token" {
  value     = zitadel_organization_domain.validated.validation_token
  sensitive = true
}

resource "zitadel_organization_domain" "verified" {
  organization_id = zitadel_organization.default.id
  domain          = "verified.example.com"
  validation_type = "DOMAIN_VALIDATION_TYPE_HTTP"
  verify          = true
}
