---
page_title: "zitadel_webkey Resource - terraform-provider-zitadel"
subcategory: ""
description: |-
	Resource representing a web key.
---

# zitadel_webkey (Resource)

Resource representing a web key.

## Example Usage

```terraform
resource "zitadel_webkey" "default" {
  org_id = data.zitadel_org.default.id
  rsa {
    bits   = "RSA_BITS_4096"
    hasher = "RSA_HASHER_SHA256"
  }
}
```

## Schema

### Required

Exactly one of the following blocks must be provided to define the key type:

- either
  - `ecdsa` (Block) Create an ECDSA key pair and specify the curve. If no curve is provided, a ECDSA key pair with P-256 curve will be created.
      - `bits` (String) Bit size of the RSA key. Default is 2048 bits.
      - `hasher` (String) Signing algrithm used. Default is SHA256.
  - `ed25519` (Block) Create a ED25519 key pair. (see [below for nested schema](#nestedblock--ed25519))
  - `rsa` (Block) Create a RSA key pair and specify the bit size and hashing algorithm. If no bits and hasher are provided, a RSA key pair with 2048 bits and SHA256 hashing will be created.
      - `curve` (String) Curve of the ECDSA key. Default is P-256.

### Optional

- `org_id` (String) ID of the organization

### Read-Only

- `id` (String) The ID of this resource.
- `key_type` (String) Type of the key.
- `state` (String) State of the key.

## Import

```bash
# The resource can be imported using the ID format `<id[:org_id]>`, e.g.
terraform import zitadel_webkey.imported '123456789012345678:123456789012345678'
```
