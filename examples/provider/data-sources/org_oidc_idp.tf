data "zitadel_org_oidc_idp" "org_oidc_idp" {
  id     = "177073612581240835"
  org_id = data.zitadel_org.org.id
}

output "org_oidc_idp" {
  value = data.zitadel_org_oidc_idp.org_oidc_idp
}
