data "zitadel_human_user" "human_user" {
  id     = "177073614158299139"
  org_id = data.zitadel_org.org.id
}

output "human_user" {
  value = data.zitadel_human_user.human_user
}
