---
page_title: "Migrating to the v2 resources"
subcategory: ""
description: |-
  How to move existing zitadel_project and zitadel_application_* resources to the
  zitadel_project_v2 and zitadel_application_v2 resources without downtime or
  losing generated client secrets.
---

# Migrating to the v2 resources

The `zitadel_project_v2` and `zitadel_application_v2` resources call ZITADEL's
v2 APIs. They are added alongside the existing `zitadel_project`,
`zitadel_application_oidc`, `zitadel_application_saml` and
`zitadel_application_api` resources, which continue to use the v1 management
API and are unchanged.

Upgrading the provider is therefore **not a breaking change**: existing
configurations keep working, and migrating to the v2 resources is entirely
opt-in. ZITADEL 3.x users must stay on the v1 resources, because the v2 APIs
only exist on ZITADEL 4.x.

## Migrate in place, without recreating anything

The v1 and v2 resources are different resource *types*. A `moved` block can
relocate state between addresses, but it copies the old state verbatim and
cannot reshape it to the v2 schema — it would not move attributes into the new
nested blocks (`oidc`/`saml`/`api`) and would drop the practitioner-supplied
`client_secret`, which the v2 API never returns again. Simply changing the
resource type in your configuration is worse still: Terraform **destroys and
recreates** the object, deleting the project or application in ZITADEL and
assigning a new ID.

Migrate with `terraform state rm` followed by `terraform import` instead: this
lets the v2 resource read the object back through the v2 API and populate the
new schema correctly. The underlying ZITADEL object IDs are identical across
the v1 and v2 APIs, so the object is never recreated — only the Terraform state
record is replaced.

### Projects

```hcl
# Before
resource "zitadel_project" "default" {
  org_id = data.zitadel_org.default.id
  name   = "my-project"
}

# After
resource "zitadel_project_v2" "default" {
  org_id = data.zitadel_org.default.id # now required
  name   = "my-project"
}
```

```sh
# read the existing id (optional, for reference)
terraform state show zitadel_project.default

# drop the v1 state entry (the project stays in ZITADEL)
terraform state rm zitadel_project.default

# import the same id into the v2 resource: <project_id[:org_id]>
terraform import zitadel_project_v2.default '123456789012345678:876543210987654321'
```

Run `terraform plan` afterwards; it should report no changes.

### Applications

`zitadel_application_oidc`, `zitadel_application_saml` and
`zitadel_application_api` all migrate to the single unified
`zitadel_application_v2` resource. The per-type configuration moves into a
nested `oidc`, `saml` or `api` block:

```hcl
# Before
resource "zitadel_application_oidc" "web" {
  project_id     = zitadel_project.default.id
  org_id         = data.zitadel_org.default.id
  name           = "web"
  redirect_uris  = ["https://localhost.com/callback"]
  response_types = ["OIDC_RESPONSE_TYPE_CODE"]
  grant_types    = ["OIDC_GRANT_TYPE_AUTHORIZATION_CODE"]
}

# After
resource "zitadel_application_v2" "web" {
  project_id = zitadel_project_v2.default.id
  org_id     = data.zitadel_org.default.id
  name       = "web"
  oidc {
    redirect_uris  = ["https://localhost.com/callback"]
    response_types = ["OIDC_RESPONSE_TYPE_CODE"]
    grant_types    = ["OIDC_GRANT_TYPE_AUTHORIZATION_CODE"]
  }
}
```

## Preserving the generated client secret

OIDC applications using `OIDC_AUTH_METHOD_TYPE_BASIC` / `_POST`, and API
applications using `API_AUTH_METHOD_TYPE_BASIC`, have a ZITADEL-generated
`client_secret`. The secret is only ever returned once, at creation time, and
the read API never returns it again, so a plain import cannot recover it.

To migrate the secret, read it from the existing v1 state and pass it as the
optional final segment of the application import ID:

```sh
# read the secret currently stored in v1 state
terraform state show zitadel_application_oidc.web   # note client_secret

terraform state rm zitadel_application_oidc.web

# import format: <app_id[:org_id[:client_secret]]>
terraform import zitadel_application_v2.web '123456789012345678:876543210987654321:THE_SECRET'
```

The import detects the application type and stores the secret in the matching
`oidc`/`api` block. A subsequent `terraform plan` reports no changes, and the
secret is **not rotated**.

Notes:

- If you do not pass the secret, the application keeps working in ZITADEL (the
  secret is not changed), but it will no longer be available in Terraform state
  or as an output. If you no longer have the secret, generate a new one with
  ZITADEL's regenerate-secret endpoint and update any consumers.
- SAML applications have no client secret, so they import with just
  `<app_id[:org_id]>`. Passing a secret segment for a SAML application is
  rejected with an error.
- If the application has no org-level scoping but you still need to pass a
  secret, supply an empty org segment: `<app_id>::<client_secret>`.
