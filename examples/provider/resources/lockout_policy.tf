resource "zitadel_lockout_policy" "default" {
  org_id                = data.zitadel_org.default.id
  max_password_attempts = "5"
}
