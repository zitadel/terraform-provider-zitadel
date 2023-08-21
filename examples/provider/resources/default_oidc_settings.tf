resource "zitadel_default_oidc_settings" "default" {
  access_token_lifetime         = "12h0m0s"
  id_token_lifetime             = "12h0m0s"
  refresh_token_expiration      = "720h0m0s"
  refresh_token_idle_expiration = "2160h0m0s"
}
