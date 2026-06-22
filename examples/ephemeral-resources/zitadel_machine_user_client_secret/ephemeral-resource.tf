# Generate a fresh client_id/client_secret pair for a machine (service) user
# without persisting the secret to Terraform state. Gate with count/for_each to
# rotate only on demand (every evaluation rotates the secret).
ephemeral "zitadel_machine_user_client_secret" "this" {
  user_id = zitadel_machine_user.this.id
}
