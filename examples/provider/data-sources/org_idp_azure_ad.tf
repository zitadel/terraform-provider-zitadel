data "zitadel_org_idp_azure_ad" "default" {
  org_id = data.zitadel_org.default.id
  id     = "123456789012345678"
}
