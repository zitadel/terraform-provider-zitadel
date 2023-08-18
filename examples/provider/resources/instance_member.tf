resource "zitadel_instance_member" "default" {
  user_id = zitadel_human_user.default.id
  roles   = ["IAM_OWNER"]
}
