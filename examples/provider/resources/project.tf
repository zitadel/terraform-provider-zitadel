resource "zitadel_project" "default" {
  name                     = "projectname"
  org_id                   = zitadel_org.default.id
  project_role_assertion   = true
  project_role_check       = true
  has_project_check        = true
  private_labeling_setting = "PRIVATE_LABELING_SETTING_ENFORCE_PROJECT_RESOURCE_OWNER_POLICY"
}
