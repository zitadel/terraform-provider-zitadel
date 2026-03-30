resource "zitadel_action_target_public_key" "default" {
  target_id  = zitadel_action_target.default.id
  public_key = file("path/to/public_key.pem")
}
