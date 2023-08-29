data "zitadel_human_user" "default" {
  org_id  = data.zitadel_org.default.id
  user_id = "123456789012345678"
}

output "human_user" {
  value = data.zitadel_human_user.default
}
