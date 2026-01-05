data "zitadel_organization_metadatas" "default" {
  organization_id = "123456789012345678"
}

data "zitadel_organization_metadatas" "filtered" {
  organization_id = "123456789012345678"
  key             = "example"
}

output "all_metadata" {
  value = data.zitadel_organization_metadatas.default.metadata
}
