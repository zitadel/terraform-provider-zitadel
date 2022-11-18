resource zitadel_domain_claimed_message_text domain_claimed_en {
  depends_on = [zitadel_org.org]

  org_id = zitadel_org.org.id
  language = "en"

  title = "title example"
  pre_header = "pre_header example"
  subject = "subject example"
  greeting = "greeting example"
  text = "text example"
  button_text = "button_text example"
  footer_text = "footer_text example"
}