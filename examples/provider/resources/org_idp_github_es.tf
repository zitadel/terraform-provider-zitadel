resource "zitadel_org_idp_github_es" "default" {
  org_id                 = data.zitadel_org.default.id
  name                   = "GitHub Enterprise Server"
  client_id              = "86a165..."
  client_secret          = "*****afdbac18"
  scopes                 = ["openid", "profile", "email"]
  authorization_endpoint = "https://auth.endpoint"
  token_endpoint         = "https://token.endpoint"
  user_endpoint          = "https://user.endpoint"
  is_linking_allowed     = false
  is_creation_allowed    = true
  is_auto_creation       = false
  is_auto_update         = true
  auto_linking        = "AUTO_LINKING_OPTION_USERNAME"
}
