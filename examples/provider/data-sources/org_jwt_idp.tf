data "zitadel_org_jwt_idp" "org_jwt_idp" {
  id     = "177073612581240835"
  org_id = data.zitadel_org.org.id
}
