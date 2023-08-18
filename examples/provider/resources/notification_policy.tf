resource "zitadel_notification_policy" "default" {
  org_id          = zitadel_org.default.id
  password_change = false
}
