resource "zitadel_project_grant_member" "default" {
  org_id     = data.zitadel_org.default.id
  project_id = data.zitadel_project.default.id
  grant_id   = data.zitadel_project_grant.default.id
  user_id    = data.zitadel_human_user.default.id
  roles      = ["PROJECT_GRANT_OWNER"]
}
