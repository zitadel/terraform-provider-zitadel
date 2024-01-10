data "zitadel_projects" "default" {
  org_id      = "123456789012345678"
  name        = "example-name"
  name_method = "TEXT_QUERY_METHOD_CONTAINS_IGNORE_CASE"
}

data "zitadel_project" "default" {
  for_each = toset(data.zitadel_projects.default.project_ids)
  id       = each.value
}

output "project_names" {
  value = toset([
    for project in data.zitadel_project.default : project.name
  ])
}
