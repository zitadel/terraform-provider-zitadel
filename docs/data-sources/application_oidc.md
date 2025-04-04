---
page_title: "zitadel_application_oidc Data Source - terraform-provider-zitadel"
subcategory: ""
description: |-
  Datasource representing an OIDC application belonging to a project, with all configuration possibilities.
---

# zitadel_application_oidc (Data Source)

Datasource representing an OIDC application belonging to a project, with all configuration possibilities.

## Example Usage

```terraform
data "zitadel_application_oidc" "default" {
  org_id     = data.zitadel_org.default.id
  project_id = data.zitadel_project.default.id
  app_id     = "123456789012345678"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `app_id` (String) The ID of this resource.
- `project_id` (String) ID of the project

### Optional

- `org_id` (String) ID of the organization

### Read-Only

- `access_token_role_assertion` (Boolean) Access token role assertion
- `access_token_type` (String) Access token type
- `additional_origins` (List of String) Additional origins
- `app_type` (String) App type
- `auth_method_type` (String) Auth method type
- `client_id` (String, Sensitive) Client ID
- `clock_skew` (String) Clockskew
- `dev_mode` (Boolean) Dev mode
- `grant_types` (List of String) Grant types
- `id` (String) The ID of this resource.
- `id_token_role_assertion` (Boolean) ID token role assertion
- `id_token_userinfo_assertion` (Boolean) Token userinfo assertion
- `name` (String) Name of the application
- `post_logout_redirect_uris` (List of String) Post logout redirect URIs
- `redirect_uris` (List of String) RedirectURIs
- `response_types` (List of String) Response type
- `skip_native_app_success_page` (Boolean) Skip the successful login page on native apps and directly redirect the user to the callback.
- `version` (String) Version