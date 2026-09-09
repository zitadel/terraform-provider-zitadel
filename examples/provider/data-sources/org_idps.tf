data "zitadel_org_idps" "default" {
  org_id      = data.zitadel_org.default.id
  name        = "example-name"
  name_method = "TEXT_QUERY_METHOD_CONTAINS_IGNORE_CASE"
  type        = "PROVIDER_TYPE_GOOGLE"
  owner_type  = "IDP_OWNER_TYPE_ORG"
}

data "zitadel_org_idp_google" "default" {
  for_each = toset(data.zitadel_org_idps.default.ids)
  org_id   = data.zitadel_org.default.id
  id       = each.value
}

output "org_idp_names" {
  value = toset([
    for idp in data.zitadel_org_idp_google.default : idp.name
  ])
}
