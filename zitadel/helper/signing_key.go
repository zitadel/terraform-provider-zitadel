package helper

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// SigningKeyRotationDiff marks a generated signing key as unknown when the
// configured expiration changes to a non-empty value. ZITADEL issues a new
// signing key whenever an update carries an expiration, so without this the
// plan keeps the prior key as a known value and anything that depends on it is
// planned with the stale key. Equivalent durations such as "0" and "0s" are
// not a change, matching DurationDiffSuppress on the expiration attribute.
func SigningKeyRotationDiff(d *schema.ResourceDiff, expirationVar, signingKeyVar string) error {
	oldValue, newValue := d.GetChange(expirationVar)
	if newValue.(string) == "" || DurationDiffSuppress(expirationVar, oldValue.(string), newValue.(string), nil) {
		return nil
	}
	return d.SetNewComputed(signingKeyVar)
}
