resource "zitadel_idp_github" "default" {
  name                = "GitHub"
  client_id           = "86a165..."
  client_secret       = "*****afdbac18"
  scopes = ["openid", "profile", "email"]
  is_linking_allowed  = false
  is_creation_allowed = true
  is_auto_creation    = false
  is_auto_update      = true
  auto_linking        = "AUTO_LINKING_OPTION_USERNAME"
}
