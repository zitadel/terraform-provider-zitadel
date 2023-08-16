data "zitadel_machine_user" "machine_user" {
  id     = "177073617463410691"
  org_id = data.zitadel_org.org.id
}

output "machine_user" {
  value = data.zitadel_machine_user.machine_user
}
