---
page_title: "Data Source zitadel_webkey - terraform-provider-zitadel"
subcategory: ""
description: |-
	Datasource representing a web key.
---

# Data Source (zitadel_webkey)

Datasource representing a web key.

## Schema

### Required

- `webkey_id` (String) The ID of this resource.

### Optional

- `org_id` (String) ID of the organization

### Read-Only

- `id` (String) The ID of this resource.
- `key_type` (String) Type of the key.
- `public_key` (String, Sensitive) The public key, PEM encoded.
- `state` (String) State of the key.
