resource "zitadel_project_grant" "default" {
  org_id         = zitadel_org.default.id
  project_id     = zitadel_project.default.id
  granted_org_id = zitadel_org.default.id
  role_keys      = [zitadel_project_role.default.role_key]
}
