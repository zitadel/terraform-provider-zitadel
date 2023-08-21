data "zitadel_org" "default" {
  id = "123456789012345678"
}

output "org" {
  value = data.zitadel_org.default
}
