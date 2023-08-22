data "zitadel_application_oidc" "default" {
  org_id     = data.zitadel_org.default.id
  project_id = data.zitadel_project.default.id
  app_id     = "123456789012345678"
}

output "application_oidc" {
  value = data.zitadel_application_oidc.default
}
