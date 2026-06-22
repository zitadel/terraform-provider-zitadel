# Regenerate the client secret of an OIDC/API application created with
# zitadel_application_v2 and hand it to a secret store, without the secret ever
# being written to Terraform state. Gate with count/for_each to rotate only on
# demand (every evaluation rotates the secret).
ephemeral "zitadel_application_v2_client_secret" "this" {
  project_id     = zitadel_application_v2.this.project_id
  application_id = zitadel_application_v2.this.id
}
