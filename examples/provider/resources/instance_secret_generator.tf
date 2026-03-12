resource "zitadel_instance_secret_generator" "password_reset" {
  generator_type        = "password_reset_code"
  length                = 8
  expiry                = "15m"
  include_lower_letters = true
  include_upper_letters = true
  include_digits        = true
  include_symbols       = false
}
