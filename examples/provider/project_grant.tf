
resource zitadel_project_grant project_grant {
  depends_on = [zitadel_org.org, zitadel_project.project, zitadel_org.grantedorg]

  org_id         = zitadel_org.org.id
  project_id     = zitadel_project.project.id
  granted_org_id = zitadel_org.grantedorg.id
}