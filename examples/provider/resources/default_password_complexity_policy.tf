resource zitadel_default_password_complexity_policy password_complexity_policy {
  min_length    = "8"
  has_uppercase = true
  has_lowercase = true
  has_number    = true
  has_symbol    = true
}