resource "zitadel_verify_sms_otp_message_text" "default" {
  org_id   = data.zitadel_org.default.id
  language = "en"

  text = "text example"
}
