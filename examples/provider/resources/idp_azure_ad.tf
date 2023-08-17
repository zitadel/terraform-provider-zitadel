resource "zitadel_idp_azure_ad" "azure_ad" {
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
}
