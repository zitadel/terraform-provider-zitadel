resource "zitadel_administrator" "default" {
  user_id       = zitadel_machine_user.default.id
  resource_type = "organization"
  org_id        = zitadel_organization.default.id
  roles         = ["ORG_USER_MANAGER"]
}

resource "zitadel_administrator" "project_grant" {
  user_id       = zitadel_machine_user.default.id
  resource_type = "project_grant"
  project_id    = zitadel_project_v2.default.id
  org_id        = zitadel_organization.grantee.id
  roles         = ["PROJECT_GRANT_OWNER"]
}
