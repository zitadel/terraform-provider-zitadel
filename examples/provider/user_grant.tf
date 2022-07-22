
resource zitade_user_grant user_grant{
  depends_on = ["zitadel_project.project", "zitadel_org.org", "zitadel_human_user.human_user"]

  project_id = zitadel_project.project.id
  org_id = zitadel_org.org.id
  role_keys = [""]
  user_id = zitadel_human_user.granted_human_user
}