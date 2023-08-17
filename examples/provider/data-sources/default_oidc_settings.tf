data "zitadel_default_oidc_settings" "oidc_settings" {}

output "oidc_settings" {
  value = data.zitadel_default_oidc_settings.oidc_settings
}
