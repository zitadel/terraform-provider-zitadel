data "zitadel_org" "default" {
  org_id = "123456789012345678"
}

output "org" {
  value = data.zitadel_org.default
}
