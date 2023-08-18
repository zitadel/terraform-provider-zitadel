resource "zitadel_password_complexity_policy" "default" {
  org_id        = data.zitadel_org.default.id
  min_length    = "8"
  has_uppercase = true
  has_lowercase = true
  has_number    = true
  has_symbol    = true
}
