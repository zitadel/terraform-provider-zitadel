data "zitadel_org_idp_google" "default" {
  org_id = data.zitadel_org.default.id
  id     = "123456789012345678"
}
