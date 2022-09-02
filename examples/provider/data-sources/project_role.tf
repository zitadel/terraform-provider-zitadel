data zitadel_project_role project_role {
  depends_on = [data.zitadel_org.org, data.zitadel_project.project]

  org_id = data.zitadel_org.org.id
  project_id     = data.zitadel_project.project.id
  role_key = "key"
}

output project_role {
  value = data.zitadel_project_role.project_role
}