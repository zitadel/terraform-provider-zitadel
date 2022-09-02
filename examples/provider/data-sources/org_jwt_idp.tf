data zitadel_org_jwt_idp org_jwt_idp {
  depends_on = [data.zitadel_org.org]

  org_id = data.zitadel_org.org.id
  idp_id = "177073612581240835"
}

output org_jwt_idp {
  value = data.zitadel_org_jwt_idp.org_jwt_idp
}