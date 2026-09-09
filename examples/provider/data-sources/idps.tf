data "zitadel_idps" "default" {
  name        = "example-name"
  name_method = "TEXT_QUERY_METHOD_CONTAINS_IGNORE_CASE"
  type        = "PROVIDER_TYPE_GITLAB"
}

data "zitadel_idp_gitlab" "default" {
  for_each = toset(data.zitadel_idps.default.ids)
  id       = each.value
}

output "idp_names" {
  value = toset([
    for idp in data.zitadel_idp_gitlab.default : idp.name
  ])
}
