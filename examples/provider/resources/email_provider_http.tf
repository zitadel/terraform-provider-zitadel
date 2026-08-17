resource "zitadel_email_provider_http" "default" {
  endpoint    = "https://example.com/provider"
  description = "provider description"
  set_active  = false
}
