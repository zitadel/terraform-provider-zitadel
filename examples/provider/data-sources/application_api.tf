data zitadel_application_api api_application {
  depends_on = [data.zitadel_org.org, data.zitadel_project.project]

  org_id     = data.zitadel_org.org.id
  project_id = data.zitadel_project.project.id
  app_id     = "177073625566806019"
}

output api_application {
  value = data.zitadel_application_api.api_application
}