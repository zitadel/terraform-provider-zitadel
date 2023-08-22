resource "zitadel_default_domain_policy" "default" {
  user_login_must_be_domain                   = false
  validate_org_domains                        = true
  smtp_sender_address_matches_instance_domain = true
}
