data "zitadel_organization_metadata" "default" {
  organization_id = "123456789012345678"
  key             = "example_key"
}

output "metadata_value" {
  value = data.zitadel_organization_metadata.default.value
}
