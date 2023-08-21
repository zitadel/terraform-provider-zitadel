resource "zitadel_user_grant" "default" {
  project_id = data.zitadel_project.default.id
  org_id     = data.zitadel_org.default.id
  role_keys  = ["some_role_key"]
  user_id    = data.zitadel_human_user.default.id
}
