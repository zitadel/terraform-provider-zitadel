resource "zitadel_machine_user" "default" {
  org_id          = data.zitadel_org.default.id
  user_name       = "machine@example.com"
  name            = "name"
  description     = "a machine user"
}
