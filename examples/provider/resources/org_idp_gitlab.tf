resource "zitadel_org_idp_gitlab" "default" {
  org_id              = data.zitadel_org.default.id
  name                = "GitLab"
  client_id           = "15765e..."
  client_secret       = "*****abcxyz"
  scopes              = ["openid", "profile", "email"]
  is_linking_allowed  = false
  is_creation_allowed = true
  is_auto_creation    = false
  is_auto_update      = true
}
