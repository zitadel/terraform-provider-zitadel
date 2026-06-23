# Create a new JSON key for a machine (service) user and hand the key material
# to a secret store, without it ever being written to Terraform state. Each
# evaluation adds a new key, so gate with count/for_each to add a key on demand.
ephemeral "zitadel_machine_key" "this" {
  user_id         = zitadel_machine_user.this.id
  expiration_date = "2519-04-01T08:45:00Z"
}
