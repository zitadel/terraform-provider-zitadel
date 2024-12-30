resource "zitadel_org_idp_oidc" "default" {
  org_id              = data.zitadel_org.default.id
  name                = "My Generic OIDC IDP"
  client_id           = "a_client_id"
  client_secret       = "a_client_secret"
  scopes              = ["openid", "profile", "email"]
  issuer              = "https://example.com"
  is_linking_allowed  = false
  is_creation_allowed = true
  is_auto_creation    = false
  is_auto_update      = true
  is_id_token_mapping = true
  auto_linking        = "AUTO_LINKING_OPTION_USERNAME"
}
