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
	if raw.IsNull() {
		return nil
	}
	if !raw.IsKnown() {
		// The whole config is not yet known (e.g. it derives from another
		// resource). Mark the hash as computed so an update is planned once the
		// secret becomes known.
		return d.SetNewComputed(hashVar)
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
	return d.SetNew(hashVar, hashSecret(v.AsString()))
}

// hashRounds stretches the secret hash to slow down offline brute-forcing of a
// low-entropy secret (e.g. a weak SMTP or LDAP password) recovered from state.
// The hash stays deterministic and unsalted on purpose: an unchanged secret
// must produce an unchanged hash so that rotation can be detected.
const hashRounds = 100000

// hashSecret returns a stretched, non-reversible hash of the secret. The value
// is only used to detect changes; it never needs to be reversed.
func hashSecret(secret string) string {
	sum := sha256.Sum256([]byte(secret))
	for i := 1; i < hashRounds; i++ {
		sum = sha256.Sum256(sum[:])
	}
	return hex.EncodeToString(sum[:])
}
