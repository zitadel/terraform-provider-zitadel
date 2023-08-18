resource "zitadel_machine_key" "default" {
  org_id          = zitadel_org.default.id
  user_id         = zitadel_machine_user.default.id
  key_type        = "KEY_TYPE_JSON"
  expiration_date = "2519-04-01T08:45:00Z"
}
