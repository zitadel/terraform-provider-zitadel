resource "zitadel_project_role" "default" {
  org_id       = zitadel_org.default.id
  project_id   = zitadel_project.default.id
  role_key     = "key"
  display_name = "display_name2"
  group        = "role_group"
}
