data "zitadel_org_idps" "default" {
  org_id = data.zitadel_org.default.id
}

data "zitadel_org_idps" "org_owned_google" {
  org_id     = data.zitadel_org.default.id
  type       = "PROVIDER_TYPE_GOOGLE"
  owner_type = "IDP_OWNER_TYPE_ORG"
}

output "org_google_idp_ids" {
  value = data.zitadel_org_idps.org_owned_google.idps[*].id
}
