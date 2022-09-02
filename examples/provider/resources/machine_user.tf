
resource zitadel_machine_user machine_user {
  depends_on  = [zitadel_org.org]

  org_id      = zitadel_org.org.id
  user_name   = "machine@localhost.com"
  name        = "name"
  description = "description"
}