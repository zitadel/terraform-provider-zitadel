---
page_title: "zitadel_login_texts Resource - terraform-provider-zitadel"
subcategory: ""
description: |-
  
---

# zitadel_login_texts (Resource)



## Example Usage

```terraform
resource zitadel_login_texts login_texts_en {
  org_id   = zitadel_org.org.id
  language = "en"

  email_verification_done_text = {
    cancel_button_text = "example"
    description        = "example"
    login_button_text  = "example"
    next_button_text   = "example"
    title              = "example"
  }
  email_verification_text = {
    code_label         = "example"
    description        = "example"
    next_button_text   = "example"
    resend_button_text = "example"
    title              = "example"
  }
  external_registration_user_overview_text = {
    back_button_text      = "example"
    description           = "example"
    email_label           = "example"
    firstname_label       = "example"
    language_label        = "example"
    lastname_label        = "example"
    next_button_text      = "example"
    nickname_label        = "example"
    phone_label           = "example"
    privacy_link_text     = "example"
    title                 = "example"
    tos_and_privacy_label = "example"
    tos_confirm           = "example"
    tos_confirm_and       = "example"
    tos_link_text         = "example"
    username_label        = "example"
  }
  external_user_not_found_text = {
    auto_register_button_text = "example"
    description               = "example"
    link_button_text          = "example"
    privacy_link_text         = "example"
    title                     = "example"
    tos_and_privacy_label     = "example"
    tos_confirm               = "example"
    tos_confirm_and           = "example"
    tos_link_text             = "example"
  }
  footer_text = {
    help           = "example"
    privacy_policy = "example"
    tos            = "example"
  }
  init_mfa_done_text = {
    cancel_button_text = "example"
    description        = "example"
    next_button_text   = "example"
    title              = "example"
  }
  init_mfa_otp_text = {
    cancel_button_text = "example"
    code_label         = "example"
    description        = "example"
    description_otp    = "example"
    next_button_text   = "example"
    secret_label       = "example"
    title              = "example"
  }
  init_mfa_prompt_text = {
    description      = "example"
    next_button_text = "example"
    otp_option       = "example"
    skip_button_text = "example"
    title            = "example"
    u2f_option       = "example"
  }
  init_mfa_u2f_text = {
    description                = "example"
    error_retry                = "example"
    not_supported              = "example"
    register_token_button_text = "example"
    title                      = "example"
    token_name_label           = "example"
  }
  init_password_done_text = {
    cancel_button_text = "example"
    description        = "example"
    next_button_text   = "example"
    title              = "example"
  }
  init_password_text = {
    code_label                 = "example"
    description                = "example"
    new_password_confirm_label = "example"
    new_password_label         = "example"
    next_button_text           = "example"
    resend_button_text         = "example"
    title                      = "example"
  }
  initialize_done_text = {
    cancel_button_text = "example"
    description        = "example"
    next_button_text   = "example"
    title              = "example"
  }
  initialize_user_text = {
    code_label                 = "example"
    description                = "example"
    new_password_confirm_label = "example"
    new_password_label         = "example"
    next_button_text           = "example"
    resend_button_text         = "example"
    title                      = "example"
  }
  linking_user_done_text = {
    cancel_button_text = "example"
    description        = "example"
    next_button_text   = "example"
    title              = "example"
  }
  login_text = {
    description                 = "example"
    description_linking_process = "example"
    external_user_description   = "example"
    login_name_label            = "example"
    login_name_placeholder      = "example"
    next_button_text            = "example"
    register_button_text        = "example"
    title                       = "example"
    title_linking_process       = "example"
    user_must_be_member_of_org  = "example"
    user_name_placeholder       = "example"
  }
  logout_text = {
    description       = "example"
    login_button_text = "example"
    title             = "example"
  }
  mfa_providers_text = {
    choose_other = "example"
    otp          = "example"
    u2f          = "example"
  }
  password_change_done_text = {
    description      = "example"
    next_button_text = "example"
    title            = "example"
  }
  password_change_text = {
    cancel_button_text         = "example"
    description                = "example"
    new_password_confirm_label = "example"
    new_password_label         = "example"
    next_button_text           = "example"
    old_password_label         = "example"
    title                      = "example"
  }
  password_reset_done_text = {
    description      = "example"
    next_button_text = "example"
    title            = "example"
  }
  password_text = {
    back_button_text = "example"
    confirmation     = "example"
    description      = "example"
    has_lowercase    = "example"
    has_number       = "example"
    has_symbol       = "example"
    has_uppercase    = "example"
    min_length       = "example"
    next_button_text = "example"
    password_label   = "example"
    reset_link_text  = "example"
    title            = "example"
  }
  passwordless_prompt_text = {
    description              = "example"
    description_init         = "example"
    next_button_text         = "example"
    passwordless_button_text = "example"
    skip_button_text         = "example"
    title                    = "example"
  }
  passwordless_registration_done_text = {
    cancel_button_text = "example"
    description        = "example"
    description_close  = "example"
    next_button_text   = "example"
    title              = "example"
  }
  passwordless_registration_text = {
    description                = "example"
    error_retry                = "example"
    not_supported              = "example"
    register_token_button_text = "example"
    title                      = "example"
    token_name_label           = "example"
  }
  passwordless_text = {
    description                = "example"
    error_retry                = "example"
    login_with_pw_button_text  = "example"
    not_supported              = "example"
    title                      = "example"
    validate_token_button_text = "example"
  }
  registration_option_text = {
    description                = "example"
    external_login_description = "example"
    title                      = "example"
    user_name_button_text      = "example"
  }
  registration_org_text = {
    description            = "example"
    email_label            = "example"
    firstname_label        = "example"
    lastname_label         = "example"
    orgname_label          = "example"
    password_confirm_label = "example"
    password_label         = "example"
    privacy_link_text      = "example"
    save_button_text       = "example"
    title                  = "example"
    tos_and_privacy_label  = "example"
    tos_confirm            = "example"
    tos_confirm_and        = "example"
    tos_link_text          = "example"
    username_label         = "example"
  }
  registration_user_text = {
    back_button_text         = "example"
    description              = "example"
    description_org_register = "example"
    email_label              = "example"
    firstname_label          = "example"
    gender_label             = "example"
    language_label           = "example"
    lastname_label           = "example"
    next_button_text         = "example"
    password_confirm_label   = "example"
    password_label           = "example"
    privacy_link_text        = "example"
    title                    = "example"
    tos_and_privacy_label    = "example"
    tos_confirm              = "example"
    tos_confirm_and          = "example"
    tos_link_text            = "example"
    username_label           = "example"
  }
  select_account_text = {
    description                 = "example"
    description_linking_process = "example"
    other_user                  = "example"
    session_state_active        = "example"
    session_state_inactive      = "example"
    title                       = "example"
    title_linking_process       = "example"
    user_must_be_member_of_org  = "example"
  }
  success_login_text = {
    auto_redirect_description = "example"
    next_button_text          = "example"
    redirected_description    = "example"
    title                     = "example"
  }
  username_change_done_text = {
    description      = "example"
    next_button_text = "example"
    title            = "example"
  }
  username_change_text = {
    cancel_button_text = "example"
    description        = "example"
    next_button_text   = "example"
    title              = "example"
    username_label     = "example"
  }
  verify_mfa_otp_text = {
    code_label       = "example"
    description      = "example"
    next_button_text = "example"
    title            = "example"
  }
  verify_mfa_u2f_text = {
    description         = "example"
    error_retry         = "example"
    not_supported       = "example"
    title               = "example"
    validate_token_text = "example"
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `language` (String)
- `org_id` (String)

### Optional

- `email_verification_done_text` (Attributes) (see [below for nested schema](#nestedatt--email_verification_done_text))
- `email_verification_text` (Attributes) (see [below for nested schema](#nestedatt--email_verification_text))
- `external_registration_user_overview_text` (Attributes) (see [below for nested schema](#nestedatt--external_registration_user_overview_text))
- `external_user_not_found_text` (Attributes) (see [below for nested schema](#nestedatt--external_user_not_found_text))
- `footer_text` (Attributes) (see [below for nested schema](#nestedatt--footer_text))
- `init_mfa_done_text` (Attributes) (see [below for nested schema](#nestedatt--init_mfa_done_text))
- `init_mfa_otp_text` (Attributes) (see [below for nested schema](#nestedatt--init_mfa_otp_text))
- `init_mfa_prompt_text` (Attributes) (see [below for nested schema](#nestedatt--init_mfa_prompt_text))
- `init_mfa_u2f_text` (Attributes) (see [below for nested schema](#nestedatt--init_mfa_u2f_text))
- `init_password_done_text` (Attributes) (see [below for nested schema](#nestedatt--init_password_done_text))
- `init_password_text` (Attributes) (see [below for nested schema](#nestedatt--init_password_text))
- `initialize_done_text` (Attributes) (see [below for nested schema](#nestedatt--initialize_done_text))
- `initialize_user_text` (Attributes) (see [below for nested schema](#nestedatt--initialize_user_text))
- `linking_user_done_text` (Attributes) (see [below for nested schema](#nestedatt--linking_user_done_text))
- `login_text` (Attributes) (see [below for nested schema](#nestedatt--login_text))
- `logout_text` (Attributes) (see [below for nested schema](#nestedatt--logout_text))
- `mfa_providers_text` (Attributes) (see [below for nested schema](#nestedatt--mfa_providers_text))
- `password_change_done_text` (Attributes) (see [below for nested schema](#nestedatt--password_change_done_text))
- `password_change_text` (Attributes) (see [below for nested schema](#nestedatt--password_change_text))
- `password_reset_done_text` (Attributes) (see [below for nested schema](#nestedatt--password_reset_done_text))
- `password_text` (Attributes) (see [below for nested schema](#nestedatt--password_text))
- `passwordless_prompt_text` (Attributes) (see [below for nested schema](#nestedatt--passwordless_prompt_text))
- `passwordless_registration_done_text` (Attributes) (see [below for nested schema](#nestedatt--passwordless_registration_done_text))
- `passwordless_registration_text` (Attributes) (see [below for nested schema](#nestedatt--passwordless_registration_text))
- `passwordless_text` (Attributes) (see [below for nested schema](#nestedatt--passwordless_text))
- `registration_option_text` (Attributes) (see [below for nested schema](#nestedatt--registration_option_text))
- `registration_org_text` (Attributes) (see [below for nested schema](#nestedatt--registration_org_text))
- `registration_user_text` (Attributes) (see [below for nested schema](#nestedatt--registration_user_text))
- `select_account_text` (Attributes) (see [below for nested schema](#nestedatt--select_account_text))
- `success_login_text` (Attributes) (see [below for nested schema](#nestedatt--success_login_text))
- `username_change_done_text` (Attributes) (see [below for nested schema](#nestedatt--username_change_done_text))
- `username_change_text` (Attributes) (see [below for nested schema](#nestedatt--username_change_text))
- `verify_mfa_otp_text` (Attributes) (see [below for nested schema](#nestedatt--verify_mfa_otp_text))
- `verify_mfa_u2f_text` (Attributes) (see [below for nested schema](#nestedatt--verify_mfa_u2f_text))

### Read-Only

- `id` (String) The ID of this resource.

<a id="nestedatt--email_verification_done_text"></a>
### Nested Schema for `email_verification_done_text`

Optional:

- `cancel_button_text` (String)
- `description` (String)
- `login_button_text` (String)
- `next_button_text` (String)
- `title` (String)


<a id="nestedatt--email_verification_text"></a>
### Nested Schema for `email_verification_text`

Optional:

- `code_label` (String)
- `description` (String)
- `next_button_text` (String)
- `resend_button_text` (String)
- `title` (String)


<a id="nestedatt--external_registration_user_overview_text"></a>
### Nested Schema for `external_registration_user_overview_text`

Optional:

- `back_button_text` (String)
- `description` (String)
- `email_label` (String)
- `firstname_label` (String)
- `language_label` (String)
- `lastname_label` (String)
- `next_button_text` (String)
- `nickname_label` (String)
- `phone_label` (String)
- `privacy_confirm` (String)
- `privacy_link_text` (String)
- `title` (String)
- `tos_and_privacy_label` (String)
- `tos_confirm` (String)
- `tos_link_text` (String)
- `username_label` (String)


<a id="nestedatt--external_user_not_found_text"></a>
### Nested Schema for `external_user_not_found_text`

Optional:

- `auto_register_button_text` (String)
- `description` (String)
- `link_button_text` (String)
- `privacy_confirm` (String)
- `privacy_link_text` (String)
- `title` (String)
- `tos_and_privacy_label` (String)
- `tos_confirm` (String)
- `tos_link_text` (String)


<a id="nestedatt--footer_text"></a>
### Nested Schema for `footer_text`

Optional:

- `help` (String)
- `privacy_policy` (String)
- `tos` (String)


<a id="nestedatt--init_mfa_done_text"></a>
### Nested Schema for `init_mfa_done_text`

Optional:

- `cancel_button_text` (String)
- `description` (String)
- `next_button_text` (String)
- `title` (String)


<a id="nestedatt--init_mfa_otp_text"></a>
### Nested Schema for `init_mfa_otp_text`

Optional:

- `cancel_button_text` (String)
- `code_label` (String)
- `description` (String)
- `description_otp` (String)
- `next_button_text` (String)
- `secret_label` (String)
- `title` (String)


<a id="nestedatt--init_mfa_prompt_text"></a>
### Nested Schema for `init_mfa_prompt_text`

Optional:

- `description` (String)
- `next_button_text` (String)
- `otp_option` (String)
- `skip_button_text` (String)
- `title` (String)
- `u2f_option` (String)


<a id="nestedatt--init_mfa_u2f_text"></a>
### Nested Schema for `init_mfa_u2f_text`

Optional:

- `description` (String)
- `error_retry` (String)
- `not_supported` (String)
- `register_token_button_text` (String)
- `title` (String)
- `token_name_label` (String)


<a id="nestedatt--init_password_done_text"></a>
### Nested Schema for `init_password_done_text`

Optional:

- `cancel_button_text` (String)
- `description` (String)
- `next_button_text` (String)
- `title` (String)


<a id="nestedatt--init_password_text"></a>
### Nested Schema for `init_password_text`

Optional:

- `code_label` (String)
- `description` (String)
- `new_password_confirm_label` (String)
- `new_password_label` (String)
- `next_button_text` (String)
- `resend_button_text` (String)
- `title` (String)


<a id="nestedatt--initialize_done_text"></a>
### Nested Schema for `initialize_done_text`

Optional:

- `cancel_button_text` (String)
- `description` (String)
- `next_button_text` (String)
- `title` (String)


<a id="nestedatt--initialize_user_text"></a>
### Nested Schema for `initialize_user_text`

Optional:

- `code_label` (String)
- `description` (String)
- `new_password_confirm_label` (String)
- `new_password_label` (String)
- `next_button_text` (String)
- `resend_button_text` (String)
- `title` (String)


<a id="nestedatt--linking_user_done_text"></a>
### Nested Schema for `linking_user_done_text`

Optional:

- `cancel_button_text` (String)
- `description` (String)
- `next_button_text` (String)
- `title` (String)


<a id="nestedatt--login_text"></a>
### Nested Schema for `login_text`

Optional:

- `description` (String)
- `description_linking_process` (String)
- `external_user_description` (String)
- `login_name_label` (String)
- `login_name_placeholder` (String)
- `next_button_text` (String)
- `register_button_text` (String)
- `title` (String)
- `title_linking_process` (String)
- `user_must_be_member_of_org` (String)
- `user_name_placeholder` (String)


<a id="nestedatt--logout_text"></a>
### Nested Schema for `logout_text`

Optional:

- `description` (String)
- `login_button_text` (String)
- `title` (String)


<a id="nestedatt--mfa_providers_text"></a>
### Nested Schema for `mfa_providers_text`

Optional:

- `choose_other` (String)
- `otp` (String)
- `u2f` (String)


<a id="nestedatt--password_change_done_text"></a>
### Nested Schema for `password_change_done_text`

Optional:

- `description` (String)
- `next_button_text` (String)
- `title` (String)


<a id="nestedatt--password_change_text"></a>
### Nested Schema for `password_change_text`

Optional:

- `cancel_button_text` (String)
- `description` (String)
- `new_password_confirm_label` (String)
- `new_password_label` (String)
- `next_button_text` (String)
- `old_password_label` (String)
- `title` (String)


<a id="nestedatt--password_reset_done_text"></a>
### Nested Schema for `password_reset_done_text`

Optional:

- `description` (String)
- `next_button_text` (String)
- `title` (String)


<a id="nestedatt--password_text"></a>
### Nested Schema for `password_text`

Optional:

- `back_button_text` (String)
- `confirmation` (String)
- `description` (String)
- `has_lowercase` (String)
- `has_number` (String)
- `has_symbol` (String)
- `has_uppercase` (String)
- `min_length` (String)
- `next_button_text` (String)
- `password_label` (String)
- `reset_link_text` (String)
- `title` (String)


<a id="nestedatt--passwordless_prompt_text"></a>
### Nested Schema for `passwordless_prompt_text`

Optional:

- `description` (String)
- `description_init` (String)
- `next_button_text` (String)
- `passwordless_button_text` (String)
- `skip_button_text` (String)
- `title` (String)


<a id="nestedatt--passwordless_registration_done_text"></a>
### Nested Schema for `passwordless_registration_done_text`

Optional:

- `cancel_button_text` (String)
- `description` (String)
- `description_close` (String)
- `next_button_text` (String)
- `title` (String)


<a id="nestedatt--passwordless_registration_text"></a>
### Nested Schema for `passwordless_registration_text`

Optional:

- `description` (String)
- `error_retry` (String)
- `not_supported` (String)
- `register_token_button_text` (String)
- `title` (String)
- `token_name_label` (String)


<a id="nestedatt--passwordless_text"></a>
### Nested Schema for `passwordless_text`

Optional:

- `description` (String)
- `error_retry` (String)
- `login_with_pw_button_text` (String)
- `not_supported` (String)
- `title` (String)
- `validate_token_button_text` (String)


<a id="nestedatt--registration_option_text"></a>
### Nested Schema for `registration_option_text`

Optional:

- `description` (String)
- `external_login_description` (String)
- `login_button_text` (String)
- `title` (String)
- `user_name_button_text` (String)


<a id="nestedatt--registration_org_text"></a>
### Nested Schema for `registration_org_text`

Optional:

- `description` (String)
- `email_label` (String)
- `firstname_label` (String)
- `lastname_label` (String)
- `orgname_label` (String)
- `password_confirm_label` (String)
- `password_label` (String)
- `privacy_confirm` (String)
- `privacy_link_text` (String)
- `save_button_text` (String)
- `title` (String)
- `tos_and_privacy_label` (String)
- `tos_confirm` (String)
- `tos_link_text` (String)
- `username_label` (String)


<a id="nestedatt--registration_user_text"></a>
### Nested Schema for `registration_user_text`

Optional:

- `back_button_text` (String)
- `description` (String)
- `description_org_register` (String)
- `email_label` (String)
- `firstname_label` (String)
- `gender_label` (String)
- `language_label` (String)
- `lastname_label` (String)
- `next_button_text` (String)
- `password_confirm_label` (String)
- `password_label` (String)
- `privacy_confirm` (String)
- `privacy_link_text` (String)
- `title` (String)
- `tos_and_privacy_label` (String)
- `tos_confirm` (String)
- `tos_link_text` (String)
- `username_label` (String)


<a id="nestedatt--select_account_text"></a>
### Nested Schema for `select_account_text`

Optional:

- `description` (String)
- `description_linking_process` (String)
- `other_user` (String)
- `session_state_active` (String)
- `session_state_inactive` (String)
- `title` (String)
- `title_linking_process` (String)
- `user_must_be_member_of_org` (String)


<a id="nestedatt--success_login_text"></a>
### Nested Schema for `success_login_text`

Optional:

- `auto_redirect_description` (String) Text to describe that auto-redirect should happen after successful login
- `next_button_text` (String)
- `redirected_description` (String) Text to describe that the window can be closed after redirect
- `title` (String)


<a id="nestedatt--username_change_done_text"></a>
### Nested Schema for `username_change_done_text`

Optional:

- `description` (String)
- `next_button_text` (String)
- `title` (String)


<a id="nestedatt--username_change_text"></a>
### Nested Schema for `username_change_text`

Optional:

- `cancel_button_text` (String)
- `description` (String)
- `next_button_text` (String)
- `title` (String)
- `username_label` (String)


<a id="nestedatt--verify_mfa_otp_text"></a>
### Nested Schema for `verify_mfa_otp_text`

Optional:

- `code_label` (String)
- `description` (String)
- `next_button_text` (String)
- `title` (String)


<a id="nestedatt--verify_mfa_u2f_text"></a>
### Nested Schema for `verify_mfa_u2f_text`

Optional:

- `description` (String)
- `error_retry` (String)
- `not_supported` (String)
- `title` (String)
- `validate_token_text` (String)