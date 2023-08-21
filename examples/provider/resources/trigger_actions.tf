resource "zitadel_trigger_actions" "default" {
  org_id       = data.zitadel_org.default.id
  flow_type    = "FLOW_TYPE_CUSTOMISE_TOKEN"
  trigger_type = "TRIGGER_TYPE_PRE_ACCESS_TOKEN_CREATION"
  action_ids   = [data.zitadel_action.default.id]
}
