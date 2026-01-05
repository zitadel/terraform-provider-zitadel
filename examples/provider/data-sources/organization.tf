data "zitadel_organization" "default" {
  id = "123456789012345678"
}

output "organization" {
  value = data.zitadel_organization.default
}
