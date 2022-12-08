resource zitadel_application_api application_api_full {
  depends_on = [zitadel_org.org, zitadel_project.project]

  org_id           = zitadel_org.org.id
  project_id       = zitadel_project.project.id
  name             = "applicationapifull"
  auth_method_type = "API_AUTH_METHOD_TYPE_PRIVATE_KEY_JWT"
}

resource zitadel_application_api application_api_min {
  depends_on = [zitadel_org.org, zitadel_project.project]

  org_id     = zitadel_org.org.id
  project_id = zitadel_project.project.id
  name       = "applicationapimin"
}