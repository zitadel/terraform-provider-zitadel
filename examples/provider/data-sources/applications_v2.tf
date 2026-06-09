data "zitadel_applications_v2" "default" {
  org_id     = data.zitadel_org.default.id
  project_id = data.zitadel_project.default.id
  name       = "example-name"
}

data "zitadel_application_v2" "default" {
  for_each   = toset(data.zitadel_applications_v2.default.app_ids)
  project_id = data.zitadel_project.default.id
  app_id     = each.value
}

output "application_v2_names" {
  value = toset([
    for app in data.zitadel_application_v2.default : app.name
  ])
}
