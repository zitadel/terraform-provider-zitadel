data "zitadel_action" "action" {
  id     = "177073621691269123"
  org_id = data.zitadel_org.org.id
}

output "action" {
  value = data.zitadel_action.action
}
