
resource zitadel_login_policy login_policy {
  depends_on = [zitadel_org.org]

  org_id                        = zitadel_org.org.id
  user_login                    = "true"
  allow_register                = "true"
  allow_external_idp            = "true"
  force_mfa                     = "false"
  passwordless_type             = "PASSWORDLESS_TYPE_ALLOWED"
  hide_password_reset           = "false"
  password_check_lifetime       = "240h"
  external_login_check_lifetime = "240h"
  multi_factor_check_lifetime   = "720h"
  mfa_init_skip_lifetime = "24h"
  second_factor_check_lifetime = "24h"
  ignore_unknown_usernames      = "true"
  default_redirect_uri          = "localhost:8080"
}