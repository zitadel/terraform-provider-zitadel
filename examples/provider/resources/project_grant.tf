resource zitadel_project_grant project_grant {
  org_id         = zitadel_org.org.id
  project_id     = zitadel_project.project.id
  granted_org_id = zitadel_org.grantedorg.id
  role_keys      = [zitadel_project_role.project_role.role_key]
}
