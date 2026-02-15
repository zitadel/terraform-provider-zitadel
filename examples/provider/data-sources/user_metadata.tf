data "zitadel_user_metadata" "default" {
  org_id  = data.zitadel_org.default.id
  user_id = "123456789012345678"
  key     = "example_key"
}

output "metadata_value" {
  value = data.zitadel_user_metadata.default.value
}
