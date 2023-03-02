resource zitadel_domain domain {
  depends_on = [zitadel_org.org]

  org_id    = zitadel_org.org.id
  name      = "zitadel.default.127.0.0.1.sslip.io"
  is_primary = true
}