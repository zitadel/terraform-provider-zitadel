resource "zitadel_sms_provider_twilio" "default" {
  sid           = "sid"
  sender_number = "019920892"
  token         = "twilio_token"
  set_active    = false
}
