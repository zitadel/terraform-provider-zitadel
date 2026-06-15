data "zitadel_project_v2" "default" {
  org_id     = data.zitadel_org.default.id
  project_id = "123456789012345678"
}
