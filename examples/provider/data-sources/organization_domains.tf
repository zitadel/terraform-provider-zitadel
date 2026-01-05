data "zitadel_organization_domains" "default" {
  organization_id = "123456789012345678"
}

data "zitadel_organization_domains" "filtered" {
  organization_id = "123456789012345678"
  domain          = "example.com"
}

output "all_domains" {
  value = data.zitadel_organization_domains.default.domains
}
