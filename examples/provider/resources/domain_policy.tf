resource zitadel_domain_policy domain_policy {
  org_id                                      = zitadel_org.org.id
  user_login_must_be_domain                   = false
  validate_org_domains                        = false
  smtp_sender_address_matches_instance_domain = false
}