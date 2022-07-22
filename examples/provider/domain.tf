
resource zitadel_domain domain {
  depends_on = [zitadel_org.org]

  org_id = zitadel_org.org.id
  name   = "localhost.com"
}