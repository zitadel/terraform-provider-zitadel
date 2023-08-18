resource "zitadel_personal_access_token" "default" {
  org_id          = data.zitadel_org.default.id
  user_id         = data.zitadel_machine_user.default.id
  expiration_date = "2519-04-01T08:45:00Z"
}
