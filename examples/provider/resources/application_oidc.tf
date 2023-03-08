resource zitadel_application_oidc application_oidc {
  project_id = zitadel_project.project.id
  org_id     = zitadel_org.org.id

  name                        = "applicationoidc"
  redirect_uris               = ["https://localhost.com"]
  response_types              = ["OIDC_RESPONSE_TYPE_CODE"]
  grant_types                 = ["OIDC_GRANT_TYPE_AUTHORIZATION_CODE"]
  post_logout_redirect_uris   = ["https://localhost.com"]
  app_type                    = "OIDC_APP_TYPE_WEB"
  auth_method_type            = "OIDC_AUTH_METHOD_TYPE_BASIC"
  version                     = "OIDC_VERSION_1_0"
  clock_skew                  = "0s"
  dev_mode                    = true
  access_token_type           = "OIDC_TOKEN_TYPE_BEARER"
  access_token_role_assertion = false
  id_token_role_assertion     = false
  id_token_userinfo_assertion = false
  additional_origins          = []
}
