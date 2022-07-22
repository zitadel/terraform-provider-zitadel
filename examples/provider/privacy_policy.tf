
resource zitadel_privacy_policy privacy_policy {
  depends_on = [zitadel_org.org]

  org_id       = zitadel_org.org.id
  tos_link     = "https://google.com"
  privacy_link = "https://google.com"
  help_link    = "https://google.com"
}