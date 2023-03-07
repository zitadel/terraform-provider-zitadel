resource zitadel_project_grant_member project_grant_member {
  org_id     = zitadel_org.org.id
  project_id = zitadel_project.project.id
  grant_id   = zitadel_project_grant.project_grant.id
  user_id    = zitadel_human_user.granted_human_user.id
  roles      = ["PROJECT_GRANT_OWNER"]
}