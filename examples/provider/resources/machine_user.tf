resource "zitadel_machine_user" "default" {
  org_id      = zitadel_org.default.id
  user_name   = "machine@localhost.com"
  name        = "name"
  description = "description"
}
