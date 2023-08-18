resource "zitadel_project_member" "default" {
  org_id     = zitadel_org.default.id
  project_id = zitadel_project.default.id
  user_id    = zitadel_human_user.default.id
  roles      = ["PROJECT_OWNER"]
}
