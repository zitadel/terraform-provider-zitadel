
resource zitadel_project project_full {
  depends_on = [zitadel_org.org]

  name                     = "projectname"
  org_id                   = zitadel_org.org.id
  project_role_assertion   = true
  project_role_check       = true
  has_project_check        = true
  private_labeling_setting = "PRIVATE_LABELING_SETTING_ENFORCE_PROJECT_RESOURCE_OWNER_POLICY"
}

resource zitadel_project project_min {
  depends_on = [zitadel_org.org]

  name                     = "projectname"
  org_id                   = zitadel_org.org.id
  private_labeling_setting = "PRIVATE_LABELING_SETTING_ENFORCE_PROJECT_RESOURCE_OWNER_POLICY"
}