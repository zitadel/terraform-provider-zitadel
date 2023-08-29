resource "zitadel_privacy_policy" "default" {
  org_id        = data.zitadel_org.default.id
  tos_link      = "https://example.com/tos"
  privacy_link  = "https://example.com/privacy"
  help_link     = "https://example.com/help"
  support_email = "support@example.com"
}
