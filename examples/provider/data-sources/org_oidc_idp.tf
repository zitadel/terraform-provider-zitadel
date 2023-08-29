data "zitadel_org_oidc_idp" "default" {
  org_id = data.zitadel_org.default.id
  id     = "123456789012345678"
}

output "org_oidc_idp" {
  value = data.zitadel_org_oidc_idp.default
}
