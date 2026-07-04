resource "zitadel_project_v2" "default" {
  name                     = "projectname"
  org_id                   = data.zitadel_org.default.id
  project_role_assertion   = true
  project_role_check       = true
  has_project_check        = true
  private_labeling_setting = "PRIVATE_LABELING_SETTING_ENFORCE_PROJECT_RESOURCE_OWNER_POLICY"
}

resource "zitadel_project_v2" "with_custom_id" {
  name       = "project-with-custom-id"
  org_id     = data.zitadel_org.default.id
  project_id = "custom-project-id-123"
}
