data "zitadel_default_oidc_settings" "default" {}

output "oidc_settings" {
  value = data.zitadel_default_oidc_settings.default
}
