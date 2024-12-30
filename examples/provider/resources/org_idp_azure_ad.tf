resource "zitadel_org_idp_azure_ad" "default" {
  org_id              = data.zitadel_org.default.id
  name                = "Azure AD"
  client_id           = "9065bfc8-a08a..."
  client_secret       = "H2n***"
  scopes              = ["openid", "profile", "email", "User.Read"]
  tenant_type         = "AZURE_AD_TENANT_TYPE_ORGANISATIONS"
  email_verified      = true
  is_linking_allowed  = false
  is_creation_allowed = true
  is_auto_creation    = false
  is_auto_update      = true
  auto_linking        = "AUTO_LINKING_OPTION_USERNAME"
}
