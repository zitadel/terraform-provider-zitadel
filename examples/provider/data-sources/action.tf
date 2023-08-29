data "zitadel_action" "default" {
  org_id    = data.zitadel_org.default.id
  action_id = "123456789012345678"
}

output "action" {
  value = data.zitadel_action.default
}
