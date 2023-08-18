resource "zitadel_domain_policy" "default" {
  org_id                                      = zitadel_org.default.id
  user_login_must_be_domain                   = false
  validate_org_domains                        = false
  smtp_sender_address_matches_instance_domain = false
}
