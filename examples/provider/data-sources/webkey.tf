data "zitadel_webkey" "default" {
  org_id    = data.zitadel_org.default.id
  webkey_id = "12345678901234DEMO"
}
