package helper

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// WriteOnlyStringValue reads the value of a write-only string attribute from the
// raw configuration. Write-only values are never stored in state, so d.Get
// always returns the zero value for them; the value is only available via the
// configuration.
func WriteOnlyStringValue(d *schema.ResourceData, attributeVar string) string {
	raw := d.GetRawConfig()
	if raw.IsNull() || !raw.IsKnown() {
		return ""
	}
	v := raw.GetAttr(attributeVar)
	if v.IsNull() || !v.IsKnown() {
		return ""
	}
	return v.AsString()
}

// WriteOnlyHashDiff keeps a computed hash attribute in sync with a write-only
// secret attribute. Because the secret itself is never stored in state,
// Terraform cannot detect when it changes; hashing the configured value into a
// tracked attribute makes a rotation show up as a normal diff and triggers an
// update, without the practitioner having to manage a version field. The hash
// is one-way and never contains the secret.
func WriteOnlyHashDiff(d *schema.ResourceDiff, secretVar, hashVar string) error {
	raw := d.GetRawConfig()
	if raw.IsNull() || !raw.IsKnown() {
		return nil
	}
	v := raw.GetAttr(secretVar)
	if v.IsNull() {
		return nil
	}
	if !v.IsKnown() {
		// The secret is not yet known (e.g. it references another resource).
		// Mark the hash as computed so a change is planned once it is known.
		return d.SetNewComputed(hashVar)
	}
	sum := sha256.Sum256([]byte(v.AsString()))
	return d.SetNew(hashVar, hex.EncodeToString(sum[:]))
}
