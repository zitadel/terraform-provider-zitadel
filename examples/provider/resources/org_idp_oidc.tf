resource "zitadel_org_idp_oidc" "default" {
  org_id               = data.zitadel_org.default.id
  name                 = "oidcidp"
  styling_type         = "STYLING_TYPE_UNSPECIFIED"
  client_id            = "a_client_id"
  client_secret        = "a_client_secret"
  issuer               = "https://google.com"
  scopes               = ["openid", "profile", "email"]
  display_name_mapping = "OIDC_MAPPING_FIELD_PREFERRED_USERNAME"
  username_mapping     = "OIDC_MAPPING_FIELD_PREFERRED_USERNAME"
  auto_register        = false
}
