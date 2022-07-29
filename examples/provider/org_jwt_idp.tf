resource zitadel_org_jwt_idp jwt_idp {
  depends_on = [zitadel_org.org]

  org_id        = zitadel_org.org.id
  name          = "jwtidp"
  styling_type  = "STYLING_TYPE_UNSPECIFIED"
  jwt_endpoint  = "https://jwtendpoint.com"
  issuer        = "https://google.com"
  keys_endpoint = "https://jwtendpoint.com/keys"
  header_name   = "x-auth-token"
  auto_register = false
}