---
page_title: "zitadel_application_saml Data Source - terraform-provider-zitadel"
subcategory: ""
description: |-
  Datasource representing a SAML application belonging to a project, with all configuration possibilities.
---

# zitadel_application_saml (Data Source)

Datasource representing a SAML application belonging to a project, with all configuration possibilities.

## Example Usage

```terraform
data "zitadel_application_saml" "default" {
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

- `id` (String) The ID of this resource.
- `metadata_xml` (String) Metadata as XML file
- `name` (String) Name of the application