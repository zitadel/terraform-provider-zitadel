resource "zitadel_application_oidc" "default" {
  project_id = data.zitadel_project.default.id
  org_id     = data.zitadel_org.default.id

  name                         = "applicationoidc"
  redirect_uris = ["https://localhost.com"]
  response_types = ["OIDC_RESPONSE_TYPE_CODE"]
  grant_types = ["OIDC_GRANT_TYPE_AUTHORIZATION_CODE"]
  post_logout_redirect_uris = ["https://localhost.com"]
  app_type                     = "OIDC_APP_TYPE_WEB"
  auth_method_type             = "OIDC_AUTH_METHOD_TYPE_BASIC"
  version                      = "OIDC_VERSION_1_0"
  clock_skew                   = "0s"
  dev_mode                     = true
  access_token_type            = "OIDC_TOKEN_TYPE_BEARER"
  access_token_role_assertion  = false
  id_token_role_assertion      = false
  id_token_userinfo_assertion  = false
  additional_origins = []
  skip_native_app_success_page = false

  # Optional: native passkey trust files (served at /.well-known/*).
  ios {
    team_id   = "ABCDE12345"
    bundle_id = "com.example.app"
  }
  android {
    package_name             = "com.example.app"
    sha256_cert_fingerprints = ["AA:BB:CC:DD:EE:FF:00:11:22:33:44:55:66:77:88:99:AA:BB:CC:DD:EE:FF:00:11:22:33:44:55:66:77:88:99"]
  }
}
