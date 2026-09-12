# validate_org_domains is an organization wide policy, so the two modes below
# are shown on separate organizations. Pick the one that matches yours.

# 1. Policy leaves org domains unvalidated (the ZITADEL default). ZITADEL
#    verifies the domain while adding it, so validation_type is not needed and
#    no challenge is generated.
resource "zitadel_domain_policy" "auto_verified" {
  org_id                                      = zitadel_organization.auto_verified.id
  user_login_must_be_domain                   = false
  validate_org_domains                        = false
  smtp_sender_address_matches_instance_domain = false
}

resource "zitadel_organization_domain" "auto_verified" {
  organization_id = zitadel_organization.auto_verified.id
  domain          = "example.com"

  depends_on = [zitadel_domain_policy.auto_verified]
}

# 2. Policy requires org domains to be validated. The domain stays unverified
#    until ownership is proven, so pick a validation_type to get a challenge,
#    publish it, and then set verify = true to have ZITADEL check it.
resource "zitadel_domain_policy" "validated" {
  org_id                                      = zitadel_organization.validated.id
  user_login_must_be_domain                   = false
  validate_org_domains                        = true
  smtp_sender_address_matches_instance_domain = false
}

resource "zitadel_organization_domain" "validated" {
  organization_id = zitadel_organization.validated.id
  domain          = "validated.example.com"
  validation_type = "DOMAIN_VALIDATION_TYPE_DNS"

  depends_on = [zitadel_domain_policy.validated]
}

output "dns_validation_token" {
  value     = zitadel_organization_domain.validated.validation_token
  sensitive = true
}

resource "zitadel_organization_domain" "verified" {
  organization_id = zitadel_organization.validated.id
  domain          = "verified.example.com"
  validation_type = "DOMAIN_VALIDATION_TYPE_HTTP"
  verify          = true

  depends_on = [zitadel_domain_policy.validated]
}
