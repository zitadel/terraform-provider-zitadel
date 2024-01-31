data "zitadel_application_saml" "default" {
  org_id     = data.zitadel_org.default.id
  project_id = data.zitadel_project.default.id
  app_id     = "123456789012345678"
}
