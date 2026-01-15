resource "zitadel_organization_domain" "default" {
  organization_id = zitadel_organization.default.id
  domain          = "example.com"
  validation_type = "DOMAIN_VALIDATION_TYPE_DNS"
}

output "dns_validation_token" {
  value     = zitadel_organization_domain.default.validation_token
  sensitive = true
}

resource "zitadel_organization_domain" "verified" {
  organization_id = zitadel_organization.default.id
  domain          = "verified.example.com"
  validation_type = "DOMAIN_VALIDATION_TYPE_HTTP"
  verify          = true
}
