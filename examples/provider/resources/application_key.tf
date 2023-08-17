resource "zitadel_application_key" "app_key" {
  org_id          = data.zitadel_org.org.id
  project_id      = data.zitadel_project.project.id
  app_id          = data.zitadel_application_api.application_api.id
  key_type        = "KEY_TYPE_JSON"
  expiration_date = "2519-04-01T08:45:00Z"
}
