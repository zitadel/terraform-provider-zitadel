resource "zitadel_org_idp_google" "default" {
  org_id              = data.zitadel_org.default.id
  name                = "Google"
  client_id           = "182902..."
  client_secret       = "GOCSPX-*****"
  scopes              = ["openid", "profile", "email"]
  is_linking_allowed  = false
  is_creation_allowed = true
  is_auto_creation    = false
  is_auto_update      = true
  auto_linking        = "AUTO_LINKING_OPTION_USERNAME"
}
