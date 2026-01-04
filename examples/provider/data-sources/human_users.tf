data "zitadel_human_users" "default" {
  org_id           = data.zitadel_org.default.id
  user_name        = "example-name"
  user_name_method = "TEXT_QUERY_METHOD_CONTAINS_IGNORE_CASE"
}

data "zitadel_human_user" "default" {
  for_each = toset(data.zitadel_human_users.default.user_ids)
  org_id  = data.zitadel_org.default.id
  user_id = each.value
}

output "user_names" {
  value = toset([
    for user in data.zitadel_human_user.default : user.user_name
  ])
}
