# (Re)generate the verification token for an organization domain without
# persisting it to Terraform state. Each evaluation generates a fresh challenge,
# so gate with count/for_each to (re)issue on demand.
ephemeral "zitadel_organization_domain_validation" "this" {
  org_id          = zitadel_organization_domain.this.organization_id
  domain          = zitadel_organization_domain.this.domain
  validation_type = "DOMAIN_VALIDATION_TYPE_DNS"
}
