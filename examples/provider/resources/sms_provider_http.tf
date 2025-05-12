resource "zitadel_sms_provider_http" "default" {
  endpoint    = "https://relay.example.com/provider"
  description = "provider description"
  set_active  = true
}
