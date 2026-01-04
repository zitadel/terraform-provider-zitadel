data "zitadel_org_idp_apple" "default" {
  id     = "123456789012345678"
  org_id = data.zitadel_org.default.id
}
