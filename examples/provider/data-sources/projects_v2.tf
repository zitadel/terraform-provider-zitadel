data "zitadel_projects_v2" "default" {
  org_id      = data.zitadel_org.default.id
  name        = "example-name"
  name_method = "TEXT_QUERY_METHOD_CONTAINS_IGNORE_CASE"
}

data "zitadel_project_v2" "default" {
  for_each   = toset(data.zitadel_projects_v2.default.project_ids)
  project_id = each.value
}

output "project_v2_names" {
  value = toset([
    for project in data.zitadel_project_v2.default : project.name
  ])
}
