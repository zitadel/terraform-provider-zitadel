resource zitadel_personal_access_token pat {
  org_id          = zitadel_org.org.id
  user_id         = zitadel_machine_user.machine_user.id
  expiration_date = "2519-04-01T08:45:00Z"
}