resource zitadel_user_grant user_grant {
  project_id = zitadel_project.project.id
  org_id     = zitadel_org.org.id
  role_keys  = ["key"]
  user_id    = zitadel_human_user.granted_human_user.id
}
