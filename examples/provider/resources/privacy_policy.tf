resource "zitadel_privacy_policy" "default" {
  org_id        = data.zitadel_org.default.id
  tos_link      = "https://google.com"
  privacy_link  = "https://google.com"
  help_link     = "https://google.com"
  support_email = "support@email.com"
}
