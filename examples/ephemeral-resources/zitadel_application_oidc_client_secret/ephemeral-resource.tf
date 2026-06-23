# Regenerate the client secret of a zitadel_application_oidc application without
# persisting it to Terraform state. Gate with count/for_each to rotate only on
# demand (every evaluation rotates the secret).
ephemeral "zitadel_application_oidc_client_secret" "this" {
  project_id = zitadel_application_oidc.this.project_id
  app_id     = zitadel_application_oidc.this.id
}
