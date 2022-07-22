
resource zitadel_project_role project_role {
  depends_on = [zitadel_org.org, zitadel_project.project]

  org_id       = zitadel_org.org.id
  project_id   = zitadel_project.project.id
  role_key     = "key"
  display_name = "display_name2"
  group        = "role_group"
}