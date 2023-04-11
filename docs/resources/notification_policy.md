---
page_title: "zitadel_notification_policy Resource - terraform-provider-zitadel"
subcategory: ""
description: |-
  Resource representing the custom notification policy of an organization.
---

# zitadel_notification_policy (Resource)

Resource representing the custom notification policy of an organization.

## Example Usage

```terraform
resource zitadel_notification_policy notification_policy {
  org_id          = zitadel_org.org.id
  password_change = false
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `org_id` (String) Id for the organization
- `password_change` (Boolean) Send notification if a user changes his password

### Read-Only

- `id` (String) The ID of this resource.