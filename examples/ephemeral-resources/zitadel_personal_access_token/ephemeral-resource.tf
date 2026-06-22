# Mint a new personal access token (PAT) for a machine (service) user and hand
# it to a secret store, without it ever being written to Terraform state. Each
# evaluation mints a new token, so gate with count/for_each to issue on demand.
ephemeral "zitadel_personal_access_token" "this" {
  user_id         = zitadel_machine_user.this.id
  expiration_date = "2519-04-01T08:45:00Z"
}
