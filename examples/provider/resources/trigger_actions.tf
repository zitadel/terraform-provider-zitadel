resource zitadel_trigger_actions trigger_actions {
  org_id       = zitadel_org.org.id
  flow_type    = "FLOW_TYPE_EXTERNAL_AUTHENTICATION"
  trigger_type = "TRIGGER_TYPE_POST_AUTHENTICATION"
  action_ids   = [zitadel_action.action.id]
}