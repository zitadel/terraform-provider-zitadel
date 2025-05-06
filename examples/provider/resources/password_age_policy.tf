resource "zitadel_password_age_policy" "default" {
  org_id           = data.zitadel_org.default.id
  max_age_days     = "30"
  expire_warn_days = "5"
}
