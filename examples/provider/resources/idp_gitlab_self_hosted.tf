resource "zitadel_idp_gitlab_self_hosted" "default" {
  name                = "GitLab Self Hosted"
  client_id           = "15765e..."
  client_secret       = "*****abcxyz"
  scopes              = ["openid", "profile", "email"]
  issuer              = "https://my.issuer"
  is_linking_allowed  = false
  is_creation_allowed = true
  is_auto_creation    = false
  is_auto_update      = true
  auto_linking        = "AUTO_LINKING_OPTION_USERNAME"
}
