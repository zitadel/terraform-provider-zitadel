
resource zitadel_lockout_policy lockout_policy {
  depends_on = [zitadel_org.org]

  org_id                = zitadel_org.org.id
  max_password_attempts = "5"
}