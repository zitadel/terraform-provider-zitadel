data "zitadel_project_roles" "default" {
  org_id     = data.zitadel_org.default.id
  project_id = data.zitadel_project.default.id
  role_key   = "admin"
}

output "project_roles" {
  value = data.zitadel_project_roles.default
}
