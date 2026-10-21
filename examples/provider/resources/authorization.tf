resource "zitadel_authorization" "default" {
  user_id    = zitadel_human_user.default.id
  project_id = zitadel_project_v2.default.id
  org_id     = zitadel_organization.default.id
  role_keys  = ["reader"]
}
