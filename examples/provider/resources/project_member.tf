resource zitadel_project_member project_member {
  org_id     = zitadel_org.org.id
  project_id = zitadel_project.project.id
  user_id    = zitadel_human_user.human_user.id
  roles      = ["PROJECT_OWNER"]
}