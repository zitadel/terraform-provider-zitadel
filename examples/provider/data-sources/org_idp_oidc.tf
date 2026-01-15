data "zitadel_org_idp_oidc" "default" {
  org_id = data.zitadel_org.default.id
  id     = "123456789012345678"
}

output "org_idp_oidc" {
  value = data.zitadel_org_idp_oidc.default
}
