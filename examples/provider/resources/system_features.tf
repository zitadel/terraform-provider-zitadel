resource "zitadel_system_features" "default" {
  login_default_org   = true
  oidc_token_exchange = true
  user_schema         = false
  improved_performance = [
    "IMPROVED_PERFORMANCE_PROJECT_GRANT",
    "IMPROVED_PERFORMANCE_ORG_DOMAIN_VERIFIED"
  ]
  oidc_single_v1_session_termination = true
  enable_back_channel_logout         = true
  login_v2 = {
    required = true
    base_uri = "https://login.example.com"
  }
  permission_check_v2 = true
}
