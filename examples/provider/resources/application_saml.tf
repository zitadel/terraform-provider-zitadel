resource "zitadel_application_saml" "default" {
  org_id       = data.zitadel_org.default.id
  project_id   = data.zitadel_project.default.id
  name         = "applicationapi"
  metadata_xml = 'metadata'
}
