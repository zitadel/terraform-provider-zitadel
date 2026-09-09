data "zitadel_user_grant" "default" {
  org_id   = data.zitadel_org.default.id
  user_id  = data.zitadel_human_user.default.id
  grant_id = "123456789012345678"
}

output "user_grant" {
  value = data.zitadel_user_grant.default
}
