data zitadel_human_user human_user {
  org_id  = data.zitadel_org.org.id
  user_id = "177073614158299139"
}

output human_user {
  value = data.zitadel_human_user.human_user
}