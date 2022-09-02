
resource zitadel_human_user granted_human_user {
  depends_on = [zitadel_org.grantedorg]

  org_id             = zitadel_org.grantedorg.id
  user_name          = "human@localhost"
  first_name         = "firstname"
  last_name          = "lastname"
  nick_name          = "nickname"
  display_name       = "displayname"
  preferred_language = "de"
  gender             = "GENDER_MALE"
  phone              = "+41799999999"
  is_phone_verified  = true
  email              = "test@zitadel.com"
  is_email_verified  = false
}
