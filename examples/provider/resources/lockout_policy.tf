resource "zitadel_lockout_policy" "default" {
  org_id                = zitadel_org.default.id
  max_password_attempts = "5"
}
