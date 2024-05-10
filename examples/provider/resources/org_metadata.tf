resource "zitadel_org_metadata" "default" {
  org_id = data.zitadel_org.default.id
  key    = "a_key"
  value  = "YV92YWx1ZQ==" // or base64encode("a_value")
}
