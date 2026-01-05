data "zitadel_organization_domain" "default" {
  organization_id = "123456789012345678"
  domain          = "example.com"
}

output "domain_verified" {
  value = data.zitadel_organization_domain.default.is_verified
}
