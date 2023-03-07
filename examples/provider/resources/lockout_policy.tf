resource zitadel_lockout_policy lockout_policy {
  org_id                = zitadel_org.org.id
  max_password_attempts = "5"
}