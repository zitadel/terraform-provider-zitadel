resource zitadel_personal_access_token pat {
  depends_on = [zitadel_human_user.human_user, zitadel_org.org]

  org_id          = zitadel_org.org.id
  user_id         = zitadel_human_user.human_user.id
  expiration_date = "2519-04-01T08:45:00.000000Z"
}