resource zitadel_login_policy login_policy {
  depends_on = [zitadel_org.org, zitadel_org_idp_jwt.jwt_idp, zitadel_org_idp_oidc.oidc_idp]

  org_id                        = zitadel_org.org.id
  user_login                    = true
  allow_register                = true
  allow_external_idp            = true
  force_mfa                     = false
  passwordless_type             = "PASSWORDLESS_TYPE_ALLOWED"
  hide_password_reset           = "false"
  password_check_lifetime       = "240h0m0s"
  external_login_check_lifetime = "240h0m0s"
  multi_factor_check_lifetime   = "24h0m0s"
  mfa_init_skip_lifetime        = "720h0m0s"
  second_factor_check_lifetime  = "24h0m0s"
  ignore_unknown_usernames      = true
  default_redirect_uri          = "localhost:8080"
  second_factors                = ["SECOND_FACTOR_TYPE_OTP", "SECOND_FACTOR_TYPE_U2F"]
  multi_factors                 = ["MULTI_FACTOR_TYPE_U2F_WITH_VERIFICATION"]
  idps                          = [zitadel_org_idp_oidc.oidc_idp.id, zitadel_org_idp_jwt.jwt_idp.id]
}
