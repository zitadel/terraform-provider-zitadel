resource zitadel_machine_key machine_key {
  org_id          = zitadel_org.org.id
  user_id         = zitadel_machine_user.machine_user.id
  key_type        = "KEY_TYPE_JSON"
  expiration_date = "2519-04-01T08:45:00Z"
}
