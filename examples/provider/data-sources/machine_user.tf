data "zitadel_machine_user" "default" {
  org_id  = data.zitadel_org.default.id
  user_id = "123456789012345678"
}

output "machine_user" {
  value = data.zitadel_machine_user.default
}
