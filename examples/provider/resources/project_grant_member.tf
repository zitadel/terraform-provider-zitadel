resource "zitadel_project_grant_member" "default" {
  org_id     = data.zitadel_org.default.id
  project_id = data.zitadel_project.default.id
  user_id    = data.zitadel_human_user.default.id
  grant_id   = "123456789012345678"
  roles = ["PROJECT_GRANT_OWNER"]
}
