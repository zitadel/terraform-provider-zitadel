package helper

import (
	"encoding/json"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// NonEmptyJSONObject returns a ValidateDiagFunc that only accepts a JSON object
// with at least one key, so arrays, scalars, null and {} fail at plan time
// instead of at apply or, for {}, as a no-op that plans the same update forever.
func NonEmptyJSONObject(attr string) schema.SchemaValidateDiagFunc {
	return func(value interface{}, _ cty.Path) diag.Diagnostics {
		s, ok := value.(string)
		if !ok {
			return diag.Errorf("%s must be a string, got %T", attr, value)
		}
		var object map[string]interface{}
		if err := json.Unmarshal([]byte(s), &object); err != nil {
			return diag.Errorf("%s must be a JSON object: %v", attr, err)
		}
		if len(object) == 0 {
			return diag.Errorf("%s must be a JSON object with at least one key", attr)
		}
		return nil
	}
}
