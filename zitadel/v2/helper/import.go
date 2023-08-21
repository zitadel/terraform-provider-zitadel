package helper

import (
	"context"
	"encoding/csv"
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
	SemicolonPlaceholder       = "__SEMICOLON__"
)

func NewImportAttribute(key string, value ConvertStringFunc, optional bool) importAttribute {
	return importAttribute{key: key, value: value, optional: optional}
}

func ImportWithID(idVar string, attrs ...importAttribute) *schema.ResourceImporter {
	return ImportWithAttributes(append([]importAttribute{NewImportAttribute(idVar, ConvertID, false)}, attrs...)...)
}

func ImportWithOptionalOrg(attrs ...importAttribute) *schema.ResourceImporter {
	return ImportWithAttributes(append([]importAttribute{ImportOptionalOrgAttribute}, attrs...)...)
}

func ImportWithIDAndOptionalOrg(idVar string, attrs ...importAttribute) *schema.ResourceImporter {
	return ImportWithID(idVar, append(attrs, ImportOptionalOrgAttribute)...)
}

func ImportWithIDAndOptionalSecret(idVar, secretKey string) *schema.ResourceImporter {
	return ImportWithID(idVar, importAttribute{key: secretKey, value: ConvertNonEmpty, optional: true})
}

func ImportWithIDAndOptionalOrgAndSecretV5(idVar, secretKey string) *schema.ResourceImporter {
	return ImportWithIDAndOptionalOrg(idVar, importAttribute{key: secretKey, value: ConvertNonEmpty, optional: true})
}

func ImportWithEmptyID(attrs ...importAttribute) *schema.ResourceImporter {
	return ImportWithAttributes(append([]importAttribute{emptyIDAttribute}, attrs...)...)
}

type ConvertStringFunc func(string) (interface{}, error)

type importAttribute struct {
	key      string
	value    ConvertStringFunc
	optional bool
}

type ImportAttributes []importAttribute

// Less makes the attributes sortable by putting the optional attributes to the end
// and the org id to the beginning of the optional attributes
func (i ImportAttributes) Less(j, k int) bool {
	left := (i)[j]
	right := (i)[k]
	if !left.optional && right.optional {
		return true
	}
	if left.optional && right.optional && left.key == OrgIDVar {
		return true
	}
	return false
}

func (i ImportAttributes) Len() int { return len(i) }

func (i ImportAttributes) Swap(j, k int) { (i)[j], (i)[k] = (i)[k], (i)[j] }

func ImportWithAttributes(attrs ...importAttribute) *schema.ResourceImporter {
	return &schema.ResourceImporter{
		StateContext: func(_ context.Context, data *schema.ResourceData, i interface{}) (ret []*schema.ResourceData, err error) {
			return []*schema.ResourceData{data}, importWithAttributes(data, attrs...)
		},
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
		if i == 0 && attr.key == `""` {
			continue
		}
		if attr.optional {
			optionalKeys = append(optionalKeys, attr.key)
		} else {
			requiredKeys = append(requiredKeys, attr.key)
		}
	}
	defer func() {
		err = ImportIDValidationError(id, requiredKeys, optionalKeys, err)
	}()
	externalMinParts := len(requiredKeys)
	internalMinParts := externalMinParts
	externalMaxParts := len(requiredKeys) + len(optionalKeys)
	internalMaxParts := len(attrs)
	csvReader := csv.NewReader(strings.NewReader(id))
	csvReader.Comma = ':'
	csvReader.LazyQuotes = true
	parts, err := csvReader.Read()
	if err != nil {
		return fmt.Errorf("failed to parse id: %w", err)
	}
	// if we expect an empty id and have more than just the emptyIDAttribute, we ensure the first part is an empty key
	if len(attrs) > 1 && attrs[0].key == emptyIDAttribute.key && parts[0] != "" {
		parts = append([]string{""}, parts...)
		internalMinParts++
	}
	if len(parts) < internalMinParts || len(parts) > internalMaxParts || internalMinParts > 0 && len(id) == 0 {
		return fmt.Errorf(`expected the number of semicolon separated parts to be between %d and %d, but got %d parts: "%s"`, externalMinParts, externalMaxParts, len(parts), strings.Join(parts, `", "`))
	}
	for i, part := range parts {
		part = strings.ReplaceAll(part, SemicolonPlaceholder, `:`)
		attr := attrs[i]
		// if the id is optional and not given, we use the emptyIDAttribute
		if attr.optional && part == "" {
			attr = emptyIDAttribute
		}
		val, err := attr.value(part)
		if err != nil {
			return fmt.Errorf("invalid value for %s: %w", attr.key, err)
		}
		if i == 0 {
			state.SetId(val.(string))
			continue
		}
		if err := state.Set(attr.key, val); err != nil {
			return fmt.Errorf("failed to set %s=%v: %w", attr.key, val, err)
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
