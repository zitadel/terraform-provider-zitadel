resource zitadel_org_idp_oidc oidc_idp {
  org_id               = zitadel_org.org.id
  name                 = "oidcidp"
  styling_type         = "STYLING_TYPE_UNSPECIFIED"
  client_id            = "google"
  client_secret        = "google_secret"
  issuer               = "https://google.com"
  scopes               = ["openid", "profile", "email"]
  display_name_mapping = "OIDC_MAPPING_FIELD_PREFERRED_USERNAME"
  username_mapping     = "OIDC_MAPPING_FIELD_PREFERRED_USERNAME"
  auto_register        = false
}