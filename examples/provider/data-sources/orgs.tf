data "zitadel_orgs" "default" {
  name          = "example-name"
  name_method   = "TEXT_QUERY_METHOD_CONTAINS_IGNORE_CASE"
  domain        = "example.com"
  domain_method = "TEXT_QUERY_METHOD_CONTAINS_IGNORE_CASE"
  state         = "ORG_STATE_ACTIVE"
}

data "zitadel_org" "default" {
  for_each = toset(data.zitadel_orgs.default.ids)
  id       = each.value
}

output "org_names" {
  value = toset([
    for org in data.zitadel_org.default : org.name
  ])
}
