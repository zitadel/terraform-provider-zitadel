resource "zitadel_project_grant" "default" {
  org_id         = data.zitadel_org.default.id
  project_id     = data.zitadel_project.default.id
  granted_org_id = data.zitadel_org.granted_org.id
  role_keys = ["super-user"]
}
