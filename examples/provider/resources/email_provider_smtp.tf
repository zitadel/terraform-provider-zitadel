resource "zitadel_email_provider_smtp" "default" {
  sender_address   = "sender@example.com"
  sender_name      = "ZITADEL"
  tls              = true
  host             = "localhost:25"
  user             = "user"
  password         = "password"
  reply_to_address = "replyto@example.com"
  description      = "SMTP email provider"
  set_active       = false
}
