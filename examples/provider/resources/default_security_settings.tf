resource "zitadel_default_security_settings" "default" {
  enable_impersonation = true

  dynamic_client_registration_enabled               = true
  dynamic_client_registration_allow_unauthenticated = false
}
