package helper

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DurationDiffSuppress suppresses differences between duration strings that represent
// the same duration (e.g., "0" and "0s", "720h" and "720h0m0s").
// This prevents false drift detection when Go's time.ParseDuration().String() normalizes duration formats.
func DurationDiffSuppress(k, old, new string, d *schema.ResourceData) bool {
	oldDuration, oldErr := time.ParseDuration(old)
	newDuration, newErr := time.ParseDuration(new)

	// If either fails to parse, fall back to string comparison
	if oldErr != nil || newErr != nil {
		return old == new
	}

	// Compare the actual duration values
	return oldDuration == newDuration
}
