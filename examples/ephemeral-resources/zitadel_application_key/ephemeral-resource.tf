# Create a new JSON key for an API application and hand the key material to a
# secret store, without it ever being written to Terraform state. Each
# evaluation adds a new key, so gate with count/for_each to add a key on demand.
ephemeral "zitadel_application_key" "this" {
  project_id      = zitadel_application_api.this.project_id
  app_id          = zitadel_application_api.this.id
  expiration_date = "2519-04-01T08:45:00Z"
}
