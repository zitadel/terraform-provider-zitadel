data zitadel_project project {
  depends_on = [data.zitadel_org.org]

  org_id = data.zitadel_org.org.id
  project_id     = "177073620768522243"
}

output project {
  value = data.zitadel_project.project
}