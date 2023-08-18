resource "zitadel_org_member" "default" {
  org_id  = zitadel_org.default.id
  user_id = zitadel_human_user.default.id
  roles   = ["ORG_OWNER"]
}
