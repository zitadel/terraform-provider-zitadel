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
	ImportOptionalOrgAttribute = NewImportAttribute(OrgIDVar, ConvertID, true)
	emptyIDAttribute           = NewImportAttribute(`""`, ConvertEmpty, false)
)

func NewImportAttribute(key string, value ConvertStringFunc, optional bool) importAttribute {
	return importAttribute{Key: key, Value: value, Optional: optional}
}

func ImportWithID(idVar string, attrs ...importAttribute) schema.StateContextFunc {
	return ImportWithAttributes(append([]importAttribute{NewImportAttribute(idVar, ConvertID, false)}, attrs...)...)
}

func ImportWithOptionalOrg(attrs ...importAttribute) schema.StateContextFunc {
	return ImportWithAttributes(append([]importAttribute{ImportOptionalOrgAttribute}, attrs...)...)
}

func ImportWithIDAndOptionalOrg(idVar string, attrs ...importAttribute) schema.StateContextFunc {
	return func(ctx context.Context, data *schema.ResourceData, i interface{}) ([]*schema.ResourceData, error) {
		return ImportWithID(idVar, append([]importAttribute{ImportOptionalOrgAttribute}, attrs...)...)(ctx, data, nil)
	}
}

func ImportWithIDAndOptionalSecret(idVar, secretKey string) schema.StateContextFunc {
	return func(ctx context.Context, data *schema.ResourceData, i interface{}) ([]*schema.ResourceData, error) {
		return ImportWithID(idVar, importAttribute{Key: secretKey, Value: ConvertNonEmpty, Optional: true})(ctx, data, nil)
	}
}

func ImportWithIDAndOptionalOrgAndSecretV5(idVar, secretKey string) schema.StateContextFunc {
	return func(ctx context.Context, data *schema.ResourceData, i interface{}) ([]*schema.ResourceData, error) {
		return ImportWithIDAndOptionalOrg(idVar, importAttribute{Key: secretKey, Value: ConvertNonEmpty, Optional: true})(ctx, data, nil)
	}
}

func ImportWithEmptyID(attrs ...importAttribute) schema.StateContextFunc {
	return func(ctx context.Context, data *schema.ResourceData, i interface{}) ([]*schema.ResourceData, error) {
		return ImportWithAttributes(append([]importAttribute{emptyIDAttribute}, attrs...)...)(ctx, data, nil)
	}
}

type ConvertStringFunc func(string) (interface{}, error)

type importAttribute struct {
	Key      string
	Value    ConvertStringFunc
	Optional bool
}

type ImportAttributes []importAttribute

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

func ImportWithAttributes(attrs ...importAttribute) schema.StateContextFunc {
	return func(_ context.Context, data *schema.ResourceData, i interface{}) (ret []*schema.ResourceData, err error) {
		return []*schema.ResourceData{data}, importWithAttributes(data, attrs...)
	}
}

type importState interface {
	Id() string
	SetId(string)
	Set(string, interface{}) error
}

func importWithAttributes(state importState, attrs ...importAttribute) (err error) {
	id := state.Id()
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
	parts := strings.Split(id, ":")
	// if we expect an empty id, we fill the first part with an empty key
	if len(attrs) > 1 && attrs[0].Key == emptyIDAttribute.Key {
		parts = append([]string{""}, parts...)
	}
	minParts := len(requiredKeys)
	externalMaxParts := minParts + len(optionalKeys)
	internalMaxParts := len(attrs)
	if len(parts) < minParts || len(parts) > internalMaxParts || minParts > 0 && len(id) == 0 {
		return fmt.Errorf("expected the number of semicolon separated parts to be between %d and %d, but got parts %v", minParts, externalMaxParts, parts)
	}
	for i, part := range parts {
		attr := attrs[i]
		// if the id is optional and not given, we use the emptyIDAttribute
		if attr.Optional && part == "" {
			attr = emptyIDAttribute
		}
		val, err := attr.Value(part)
		if err != nil {
			return fmt.Errorf("invalid value for %s: %w", attr.Key, err)
		}
		if i == 0 {
			state.SetId(val.(string))
			if attr.Key == ResourceIDVar || attr.Key == emptyIDAttribute.Key {
				continue
			}
		}
		if err := state.Set(attr.Key, val); err != nil {
			return fmt.Errorf("failed to set %s=%v: %w", attr.Key, val, err)
		}
	}
	return nil
}

var _ ConvertStringFunc = ConvertID

func ConvertID(id string) (interface{}, error) {
	if !ZitadelGeneratedIdOnlyRegex.MatchString(id) {
		return nil, fmt.Errorf(`id "%s" does not match regular expression %s`, id, ZitadelGeneratedIdOnlyRegex.String())
	}
	return id, nil
}

var _ ConvertStringFunc = ConvertJSON

func ConvertJSON(importValue string) (interface{}, error) {
	if err := json.Unmarshal([]byte(importValue), &struct{}{}); err != nil {
		return nil, fmt.Errorf("value must be valid JSON: %w", err)
	}
	return importValue, nil
}

var _ ConvertStringFunc = ConvertEmpty

func ConvertEmpty(importValue string) (interface{}, error) {
	if len(importValue) > 0 {
		return nil, fmt.Errorf(`value must be empty, but got "%s"`, importValue)
	}
	return "imported", nil
}

var _ ConvertStringFunc = ConvertNonEmpty

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
	expectFormat = strings.Replace(expectFormat, "<[:", "<[", 1)
	expectFormat += ">"
	return fmt.Errorf(`failed to import id "%s" by format %s: %w`, givenID, expectFormat, err)
}
