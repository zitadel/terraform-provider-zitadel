package helper

import (
	"context"
	"encoding/base64"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"sort"
	"strconv"
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

// ImportWithID is a convenience function that calls ImportWithAttributes.
// It returns a ResourceImporter that expects a ZITADEL ID number at the first import string position along with other given attributes.
// idVar is only relevant for the error message, the resources SetID function is called with first argument ID
func ImportWithID(idVar string, attributes ...importAttribute) *schema.ResourceImporter {
	return ImportWithAttributes(append([]importAttribute{NewImportAttribute(idVar, ConvertID, false)}, attributes...)...)
}

// ImportWithOptionalOrg is a convenience function that calls ImportWithAttributes.
// It returns a ResourceImporter that accepts an optional organization id along with other given attributes
func ImportWithOptionalOrg(attributes ...importAttribute) *schema.ResourceImporter {
	return ImportWithAttributes(append([]importAttribute{ImportOptionalOrgAttribute}, attributes...)...)
}

// ImportWithIDAndOptionalOrg is a convenience function that calls ImportWithID
// and passes an optional attribute for the org ID along with the other given attributes.
func ImportWithIDAndOptionalOrg(idVar string, attributes ...importAttribute) *schema.ResourceImporter {
	return ImportWithID(idVar, append(attributes, ImportOptionalOrgAttribute)...)
}

// ImportWithIDAndOptionalSecret is a convenience function that calls ImportWithID
// and passes an optional attribute for the secret var at secretKey.
func ImportWithIDAndOptionalSecret(idVar, secretKey string) *schema.ResourceImporter {
	return ImportWithID(idVar, importAttribute{key: secretKey, value: ConvertNonEmpty, optional: true})
}

// ImportWithIDAndOptionalOrgAndSecret is a convenience function that calls ImportWithIDAndOptionalOrg
// and passes an optional attribute for the secret var at secretKey.
func ImportWithIDAndOptionalOrgAndSecret(idVar, secretKey string) *schema.ResourceImporter {
	return ImportWithIDAndOptionalOrg(idVar, importAttribute{key: secretKey, value: ConvertNonEmpty, optional: true})
}

// ImportWithEmptyID returns a ResourceImporter that does not use the first import string position value
// for the states SetID call. It uses a dummy value, instead.
func ImportWithEmptyID(attributes ...importAttribute) *schema.ResourceImporter {
	return ImportWithAttributes(append([]importAttribute{emptyIDAttribute}, attributes...)...)
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

// importWithAttributes imports a resources state that is needed to query the remote resource
// as well as state that is not readable from the ZITADEL API
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
	if err != nil && err != io.EOF {
		return fmt.Errorf("failed to parse id: %w", err)
	}
	// if we expect an empty id and have more than just the emptyIDAttribute, we prepend an empty part to the ID
	if len(attrs) > 0 && attrs[0].key == emptyIDAttribute.key || attrs[0].optional && len(parts) == 0 {
		parts = append([]string{""}, parts...)
		internalMinParts++
	}
	if len(parts) < internalMinParts || len(parts) > internalMaxParts {
		return fmt.Errorf(`expected the number of semicolon separated parts to be between %d and %d, but got %d parts: "%s"`, externalMinParts, externalMaxParts, len(parts), strings.Join(parts, `", "`))
	}
	for i, part := range parts {
		part = strings.ReplaceAll(part, SemicolonPlaceholder, ":")
		attr := attrs[i]
		// if the id is optional and not given, we use the emptyIDAttribute
		if i == 0 && attr.optional && part == "" {
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
		if val == nil && attr.optional {
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
	if len(importValue) == 0 {
		return nil, nil
	}
	if err := json.Unmarshal([]byte(importValue), &struct{}{}); err != nil {
		return nil, fmt.Errorf("value must be valid JSON: %w", err)
	}
	return importValue, nil
}

func ConvertBase64(importValue string) (interface{}, error) {
	importValueDecoded, err := base64.StdEncoding.DecodeString(importValue)
	if err != nil {
		return nil, fmt.Errorf("value must be valid base64: %w", err)
	}
	return string(importValueDecoded), nil
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

func ConvertBool(importValue string) (interface{}, error) {
	return strconv.ParseBool(importValue)
}

// ImportIDValidationError wraps err with a help message about the expected format if it is not nil
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
