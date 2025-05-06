resource "zitadel_default_password_age_policy" "default" {
  max_age_days     = "30"
  expire_warn_days = "5"
}
