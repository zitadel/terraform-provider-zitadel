data "zitadel_user_metadatas" "default" {
  org_id  = data.zitadel_org.default.id
  user_id = "123456789012345678"
}

data "zitadel_user_metadatas" "filtered" {
  org_id  = data.zitadel_org.default.id
  user_id = "123456789012345678"
  key     = "example"
}

output "all_metadata" {
  value = data.zitadel_user_metadatas.default.metadata
}
