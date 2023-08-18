resource "zitadel_domain" "default" {
  org_id     = zitadel_org.default.id
  name       = "zitadel.default.127.0.0.1.sslip.io"
  is_primary = true
}
