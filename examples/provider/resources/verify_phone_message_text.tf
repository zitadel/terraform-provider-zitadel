resource "zitadel_verify_phone_message_text" "default" {
  org_id   = data.zitadel_org.default.id
  language = "en"

  text = "text example"
}
