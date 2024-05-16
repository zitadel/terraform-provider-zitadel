resource "zitadel_user_metadata" "default" {
  org_id  = data.zitadel_org.default.id
  user_id = data.zitadel_human_user.default.id
  key     = "a_key"
  value   = "a_value"
}
