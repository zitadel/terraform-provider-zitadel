data "zitadel_application_oidc" "oidc_application" {
  id         = "177073626925760515"
  org_id     = data.zitadel_org.org.id
  project_id = data.zitadel_project.project.id
}

output "oidc_application" {
  value = data.zitadel_application_oidc.oidc_application
}
