data "zitadel_application_apis" "default" {
  org_id      = "123456789012345678"
  project_id = "234567890123456789"
  name        = "example-name"
  name_method = "TEXT_QUERY_METHOD_CONTAINS_IGNORE_CASE"
}

data "zitadel_application_api" "default" {
  for_each = toset(data.zitadel_application_apis.default.app_ids)
  id       = each.value
}

output "app_api_names" {
  value = toset([
    for app in data.zitadel_application_api.default : app.name
  ])
}
