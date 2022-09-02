data zitadel_action action {
  depends_on = [data.zitadel_org.org]

  org_id     = data.zitadel_org.org.id
  action_id         = "177073621691269123"
}

output action {
  value = data.zitadel_action.action
}