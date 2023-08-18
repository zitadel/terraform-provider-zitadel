resource "zitadel_org_idp_jwt" "default" {
  org_id        = zitadel_org.default.id
  name          = "jwtidp"
  styling_type  = "STYLING_TYPE_UNSPECIFIED"
  jwt_endpoint  = "https://jwtendpoint.com"
  issuer        = "https://google.com"
  keys_endpoint = "https://jwtendpoint.com/keys"
  header_name   = "x-auth-token"
  auto_register = false
}
