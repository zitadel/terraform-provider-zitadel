resource "zitadel_project_grant_member" "default" {
  org_id     = zitadel_org.default.id
  project_id = zitadel_project.default.id
  grant_id   = zitadel_project_grant.default.id
  user_id    = zitadel_human_user.default.id
  roles      = ["PROJECT_GRANT_OWNER"]
}
