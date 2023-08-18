resource "zitadel_application_key" "default" {
  org_id          = data.zitadel_org.default.id
  project_id      = data.zitadel_project.default.id
  app_id          = data.zitadel_application_api.default.id
  key_type        = "KEY_TYPE_JSON"
  expiration_date = "2519-04-01T08:45:00Z"
}
