resource "zitadel_project_member" "default" {
  org_id     = data.zitadel_org.default.id
  project_id = data.zitadel_project.default.id
  user_id    = data.zitadel_human_user.default.id
  roles      = ["PROJECT_OWNER"]
}
