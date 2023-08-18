resource "zitadel_password_reset_message_text" "default" {
  org_id   = data.zitadel_org.default.id
  language = "en"

  title       = "title example"
  pre_header  = "pre_header example"
  subject     = "subject example"
  greeting    = "greeting example"
  text        = "text example"
  button_text = "button_text example"
  footer_text = "footer_text example"
}
