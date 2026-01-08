data "zitadel_organizations" "default" {
  name          = "example-name"
  name_method   = "TEXT_QUERY_METHOD_CONTAINS_IGNORE_CASE"
  domain        = "example.com"
  domain_method = "TEXT_QUERY_METHOD_CONTAINS_IGNORE_CASE"
  state         = "ORGANIZATION_STATE_ACTIVE"
}

data "zitadel_organizations" "default_org" {
  is_default = true
}

data "zitadel_organization" "default" {
  for_each = toset(data.zitadel_organizations.default.ids)
  id = each.value
}

output "organization_names" {
  value = toset([
    for org in data.zitadel_organization.default : org.name
  ])
}
