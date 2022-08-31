resource zitadel_application_key app_key {
  depends_on = [zitadel_application_api.application_api, zitadel_project.project, zitadel_org.org]

  org_id          = zitadel_org.org.id
  project_id      = zitadel_project.project.id
  app_id          = zitadel_application_api.application_api.id
  key_type        = "KEY_TYPE_JSON"
  expiration_date = "2519-04-01T08:45:00Z"
}