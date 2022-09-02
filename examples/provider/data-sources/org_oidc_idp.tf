data zitadel_org_oidc_idp org_oidc_idp {
  depends_on = [data.zitadel_org.org]

  org_id       = data.zitadel_org.org.id
  idp_id = "177073612581240835"
}

output org_oidc_idp {
  value = data.zitadel_org_oidc_idp.org_oidc_idp
}