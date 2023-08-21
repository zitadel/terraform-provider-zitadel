resource "zitadel_application_api" "default" {
  org_id           = data.zitadel_org.default.id
  project_id       = data.zitadel_project.default.id
  name             = "applicationapi"
  auth_method_type = "API_AUTH_METHOD_TYPE_PRIVATE_KEY_JWT"
}
