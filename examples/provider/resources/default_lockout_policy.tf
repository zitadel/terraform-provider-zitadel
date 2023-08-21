resource "zitadel_default_lockout_policy" "default" {
  max_password_attempts = "5"
}
