resource "zitadel_smtp_config" "default" {
  sender_address = "sender@example.com"
  sender_name    = "no-reply"
  tls            = true
  host           = "localhost:25"
  user           = "user"
  password       = "secret_password"
}
