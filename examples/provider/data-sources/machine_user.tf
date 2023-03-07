data zitadel_machine_user machine_user {
  org_id  = data.zitadel_org.org.id
  user_id = "177073617463410691"
}

output machine_user {
  value = data.zitadel_machine_user.machine_user
}