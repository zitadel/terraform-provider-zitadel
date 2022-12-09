resource zitadel_human_user human_user {
  depends_on = [zitadel_org.org]

  org_id             = zitadel_org.org.id
  user_name          = "humanfull@localhost.com"
  first_name         = "firstname"
  last_name          = "lastname"
  nick_name          = "nickname"
  display_name       = "displayname"
  preferred_language = "de"
  gender             = "GENDER_MALE"
  phone              = "+41799999999"
  is_phone_verified  = true
  email              = "test@zitadel.com"
  is_email_verified  = true
  initial_password   = "Password1!"
}
