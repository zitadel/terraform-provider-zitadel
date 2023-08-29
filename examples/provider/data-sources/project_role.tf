data "zitadel_project_role" "default" {
  org_id     = data.zitadel_org.default.id
  project_id = data.zitadel_project.default.id
  role_key   = "key"
}

output "project_role" {
  value = data.zitadel_project_role.default
}
