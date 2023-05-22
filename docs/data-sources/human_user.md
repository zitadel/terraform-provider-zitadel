---
page_title: "zitadel_human_user Data Source - terraform-provider-zitadel"
subcategory: ""
description: |-
  Datasource representing a human user situated under an organization, which then can be authorized through memberships or direct grants on other resources.
---

# zitadel_human_user (Data Source)

Datasource representing a human user situated under an organization, which then can be authorized through memberships or direct grants on other resources.

## Example Usage

```terraform
data zitadel_human_user human_user {
  org_id  = data.zitadel_org.org.id
  user_id = "177073614158299139"
}

output human_user {
  value = data.zitadel_human_user.human_user
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `org_id` (String) ID of the organization
- `user_id` (String) The ID of this resource.

### Read-Only

- `display_name` (String) Display name of the user
- `email` (String) Email of the user
- `first_name` (String) First name of the user
- `gender` (String) Gender of the user
- `id` (String) The ID of this resource.
- `is_email_verified` (Boolean) Is the email verified of the user, can only be true if password of the user is set
- `is_phone_verified` (Boolean) Is the phone verified of the user
- `last_name` (String) Last name of the user
- `login_names` (List of String) Loginnames
- `nick_name` (String) Nick name of the user
- `phone` (String) Phone of the user
- `preferred_language` (String) Preferred language of the user
- `preferred_login_name` (String) Preferred login name
- `state` (String) State of the user
- `user_name` (String) Username