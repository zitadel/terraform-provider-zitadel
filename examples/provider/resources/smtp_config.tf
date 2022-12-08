resource zitadel_smtp_config smtp_full {
  sender_address = "address"
  sender_name    = "no-reply"
  tls            = true
  host           = "localhost:25"
  user           = "user"
  password       = "password"
}

resource zitadel_smtp_config smtp_min {
  sender_address = "address"
  sender_name    = "no-reply"
  host           = "localhost:25"
}