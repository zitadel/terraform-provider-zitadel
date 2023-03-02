resource zitadel_notification_policy notification_policy {
  depends_on = [zitadel_org.org]

  org_id          = zitadel_org.org.id
  password_change = false
}