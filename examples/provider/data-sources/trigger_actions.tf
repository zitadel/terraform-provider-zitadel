data "zitadel_trigger_actions" "default" {
  org_id       = data.zitadel_org.default.id
  flow_type    = "FLOW_TYPE_EXTERNAL_AUTHENTICATION"
  trigger_type = "TRIGGER_TYPE_POST_AUTHENTICATION"
}

output "trigger_actions" {
  value = data.zitadel_trigger_actions.default
}
