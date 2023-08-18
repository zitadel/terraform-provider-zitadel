resource "zitadel_user_grant" "default" {
  project_id = zitadel_project.default.id
  org_id     = zitadel_org.default.id
  role_keys  = ["key"]
  user_id    = zitadel_human_user.default.id
}
