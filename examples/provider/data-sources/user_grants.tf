data "zitadel_user_grants" "default" {
  org_id     = data.zitadel_org.default.id
  user_id    = data.zitadel_human_user.default.id
  project_id = data.zitadel_project.default.id
  role_key   = "example-role"
}

data "zitadel_user_grant" "default" {
  for_each = toset(data.zitadel_user_grants.default.ids)
  org_id   = data.zitadel_org.default.id
  user_id  = data.zitadel_human_user.default.id
  grant_id = each.value
}

output "user_grant_project_ids" {
  value = toset([
    for grant in data.zitadel_user_grant.default : grant.project_id
  ])
}
