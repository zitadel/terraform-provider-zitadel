resource "zitadel_org_member" "default" {
  org_id  = data.zitadel_org.default.id
  user_id = data.zitadel_human_user.default.id
  roles = ["ORG_OWNER"]
}
