
resource zitadel_password_complexity_policy password_complexity_policy {
  depends_on = [zitadel_org.org]

  org_id        = zitadel_org.org.id
  min_length    = "8"
  has_uppercase = true
  has_lowercase = true
  has_number    = true
  has_symbol    = true
}