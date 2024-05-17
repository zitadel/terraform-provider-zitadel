resource "zitadel_idp_oauth" "default" {
  name                   = "GitLab"
  client_id              = "15765e..."
  client_secret          = "*****abcxyz"
  authorization_endpoint = "https://accounts.google.com/o/oauth2/v2/auth"
  token_endpoint         = "https://oauth2.googleapis.com/token"
  user_endpoint          = "https://openidconnect.googleapis.com/v1/userinfo"
  id_attribute           = "user_id"
  scopes                 = ["openid", "profile", "email"]
  is_linking_allowed     = false
  is_creation_allowed    = true
  is_auto_creation       = false
  is_auto_update         = true
}
