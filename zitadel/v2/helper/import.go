package helper

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	ImportOptionalOrgAttribute = ImportAttribute{Key: OrgIDVar, ValueFromString: ConvertID, Optional: true}
)

func ImportWithID(idVar string, attrs ...ImportAttribute) schema.StateContextFunc {
	return ImportWithAttributesV5(append([]ImportAttribute{{Key: idVar, ValueFromString: ConvertID}}, attrs...)...)
}

func ImportWithOptionalOrg(attrs ...ImportAttribute) schema.StateContextFunc {
	return ImportWithAttributesV5(append([]ImportAttribute{ImportOptionalOrgAttribute}, attrs...)...)
}

func ImportWithIDAndOptionalOrg(idVar string, attrs ...ImportAttribute) schema.StateContextFunc {
	return func(ctx context.Context, data *schema.ResourceData, i interface{}) ([]*schema.ResourceData, error) {
		return ImportWithID(idVar, append([]ImportAttribute{ImportOptionalOrgAttribute}, attrs...)...)(ctx, data, nil)
	}
}

func ImportWithIDAndOptionalSecret(idVar, secretKey string) schema.StateContextFunc {
	return func(ctx context.Context, data *schema.ResourceData, i interface{}) ([]*schema.ResourceData, error) {
		return ImportWithID(idVar, ImportAttribute{Key: secretKey, ValueFromString: ConvertNonEmpty, Optional: true})(ctx, data, nil)
	}
}

func ImportWithIDAndOptionalOrgAndSecretV5(idVar, secretKey string) schema.StateContextFunc {
	return func(ctx context.Context, data *schema.ResourceData, i interface{}) ([]*schema.ResourceData, error) {
		return ImportWithIDAndOptionalOrg(idVar, ImportAttribute{Key: secretKey, ValueFromString: ConvertNonEmpty, Optional: true})(ctx, data, nil)
	}
}

func ImportWithEmptyIDV5(attrs ...ImportAttribute) schema.StateContextFunc {
	return func(ctx context.Context, data *schema.ResourceData, i interface{}) ([]*schema.ResourceData, error) {
		return ImportWithAttributesV5(append([]ImportAttribute{{
			Key:             `""`,
			ValueFromString: ConvertEmpty,
		}}, attrs...)...)(ctx, data, nil)
	}
}

type ImportAttribute struct {
	Key             string
	ValueFromString func(string) (interface{}, error)
	Optional        bool
}

type ImportAttributes []ImportAttribute

// Less makes the attributes sortable by putting the optional attributes to the end
// and the org id to the beginning of the optional attributes
func (i ImportAttributes) Less(j, k int) bool {
	left := (i)[j]
	right := (i)[k]
	if !left.Optional && right.Optional {
		return true
	}
	if left.Optional && right.Optional && left.Key == OrgIDVar {
		return true
	}
	return false
}

func (i ImportAttributes) Len() int { return len(i) }

func (i ImportAttributes) Swap(j, k int) { (i)[j], (i)[k] = (i)[k], (i)[j] }

func ImportWithAttributesV5(attrs ...ImportAttribute) schema.StateContextFunc {
	return func(ctx context.Context, data *schema.ResourceData, i interface{}) (ret []*schema.ResourceData, err error) {
		id := data.Id()
		var (
			optionalKeys []string
			requiredKeys []string
		)
		sort.Sort(ImportAttributes(attrs))
		for i, attr := range attrs {
			if i == 0 && attr.Key == `""` {
				continue
			}
			if attr.Optional {
				optionalKeys = append(optionalKeys, attr.Key)
			} else {
				requiredKeys = append(requiredKeys, attr.Key)
			}
		}
		defer func() {
			err = ImportIDValidationError(id, requiredKeys, optionalKeys, err)
		}()
		parts := strings.SplitN(id, ":", len(attrs))
		// if the id should be empty and we expect an empty id, we fill the first part with an empty key
		if len(attrs) > 1 && attrs[0].Key == `""` {
			parts = append([]string{""}, parts...)
		}
		minParts := len(requiredKeys)
		maxParts := len(attrs)
		if len(parts) < minParts || len(parts) > maxParts {
			return nil, fmt.Errorf("expected the number of semicolon separated parts to be within %d and %d, but got %s", minParts, maxParts, (parts))
		}
		for i, part := range parts {
			attr := attrs[i]
			val, err := attr.ValueFromString(part)
			if err != nil {
				return nil, fmt.Errorf("invalid value %s for %s: %w", part, attr.Key, err)
			}
			if i == 0 {
				data.SetId(val.(string))
				if attr.Key == ResourceIDVar || attr.Key == `""` {
					continue
				}
			}
			if err := data.Set(attr.Key, val); err != nil {
				return nil, fmt.Errorf("failed to set %s=%v: %w", attr.Key, val, err)
			}
		}
		return []*schema.ResourceData{data}, nil
	}
}

func ConvertID(id string) (interface{}, error) {
	if !ZitadelGeneratedIdOnlyRegex.MatchString(id) {
		return nil, fmt.Errorf("id does not match regular expression %s", ZitadelGeneratedIdOnlyRegex.String())
	}
	return id, nil
}

func ConvertJSON(importValue string) (interface{}, error) {
	if err := json.Unmarshal([]byte(importValue), &struct{}{}); err != nil {
		return nil, fmt.Errorf("value must be valid JSON: %w", err)
	}
	return importValue, nil
}

func ConvertEmpty(importValue string) (interface{}, error) {
	if len(importValue) > 0 {
		return nil, fmt.Errorf("value must be empty, but got %s", importValue)
	}
	return "imported", nil
}

func ConvertNonEmpty(importValue string) (interface{}, error) {
	if len(importValue) == 0 {
		return nil, errors.New("value must not be empty")
	}
	return importValue, nil
}

func ImportIDValidationError(givenID string, requiredKeys, optionalKeys []string, err error) error {
	if err == nil {
		return nil
	}
	expectFormat := fmt.Sprintf("<")
	if len(requiredKeys) > 0 {
		expectFormat += fmt.Sprintf(":%s", strings.Join(requiredKeys, ":"))
	}
	if len(optionalKeys) > 0 {
		expectFormat += fmt.Sprintf("[:%s]", strings.Join(optionalKeys, "][:"))
	}
	expectFormat = strings.Replace(expectFormat, "<:", "<", 1)
	expectFormat += ">"
	return fmt.Errorf("failed to import id %s by format %s: %w", givenID, expectFormat, err)
}
