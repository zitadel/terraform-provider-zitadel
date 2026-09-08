data "zitadel_user_grants" "default" {
  org_id  = data.zitadel_org.default.id
  user_id = "123456789012345678"
}

data "zitadel_user_grants" "filtered" {
  org_id     = data.zitadel_org.default.id
  user_id    = "123456789012345678"
  project_id = data.zitadel_project.default.id
}

output "all_user_grants" {
  value = data.zitadel_user_grants.default.user_grants
}
