data "zitadel_project" "project" {
  id     = "177073620768522243"
  org_id = data.zitadel_org.org.id
}

output "project" {
  value = data.zitadel_project.project
}
