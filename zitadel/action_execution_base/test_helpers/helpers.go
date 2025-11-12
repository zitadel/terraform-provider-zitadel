package test_helpers

import (
	"slices"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

func CheckSchemaExactlyOneOfConsistency(t *testing.T, resourceSchema map[string]*schema.Schema) {
	t.Helper()

	if resourceSchema == nil {
		t.Fatal("resource schema is nil")
	}

	validatedGroups := make(map[string]bool)

	for sourceKey, s := range resourceSchema {
		if s.ExactlyOneOf == nil || len(s.ExactlyOneOf) == 0 {
			continue
		}

		sourceGroup := s.ExactlyOneOf
		sortedSourceGroup := slices.Clone(sourceGroup)
		slices.Sort(sortedSourceGroup)

		groupKey := strings.Join(sortedSourceGroup, ",")

		if validatedGroups[groupKey] {
			continue
		}

		for _, memberKey := range sortedSourceGroup {
			memberSchema, ok := resourceSchema[memberKey]
			if !ok {
				t.Errorf("schema key %q ExactlyOneOf list references non-existent key %q", sourceKey, memberKey)
				continue
			}

			if memberSchema.ExactlyOneOf == nil {
				t.Errorf("schema key %q is in an ExactlyOneOf group, but its member %q does not have ExactlyOneOf set", sourceKey, memberKey)
				continue
			}

			assert.ElementsMatch(t, sortedSourceGroup, memberSchema.ExactlyOneOf,
				"schema consistency error: Key %q Group %v and Member Key %q Group %v do not match",
				sourceKey, sortedSourceGroup, memberKey, memberSchema.ExactlyOneOf)
		}

		validatedGroups[groupKey] = true
	}
}

func CheckSchemaInternalValidation(t *testing.T, resource *schema.Resource, writable bool) {
	t.Helper()

	if resource == nil {
		t.Fatal("resource is nil")
	}

	if err := resource.InternalValidate(nil, writable); err != nil {
		t.Fatalf("schema.InternalValidate() failed: %s", err)
	}
}
