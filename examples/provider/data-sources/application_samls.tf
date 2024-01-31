data "zitadel_application_samls" "default" {
  org_id      = data.zitadel_org.default.id
  project_id  = data.zitadel_project.default.id
  name        = "example-name"
  name_method = "TEXT_QUERY_METHOD_CONTAINS_IGNORE_CASE"
}

data "zitadel_application_saml" "default" {
  for_each = toset(data.zitadel_application_samls.default.app_ids)
  id       = each.value
}

output "app_saml_names" {
  value = toset([
    for app in data.zitadel_application_saml.default : app.name
  ])
}
