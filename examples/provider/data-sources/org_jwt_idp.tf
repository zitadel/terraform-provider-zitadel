data "zitadel_org_jwt_idp" "default" {
  org_id = data.zitadel_org.default.id
  id     = "123456789012345678"
}

output "org_idp_org_jwt_idp" {
  value = data.zitadel_org_jwt_idp.default
}
