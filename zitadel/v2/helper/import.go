package helper

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	_                  schema.StateContextFunc = ImportWithIDV5
	_                  schema.StateContextFunc = ImportWithIDAndOrgV5
	ImportOrgAttribute                         = ImportAttribute{Key: OrgIDVar, ValueFromString: ConvertID}
)

func ImportWithIDV5(_ context.Context, data *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	if _, err := ConvertID(data.Id()); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{data}, nil
}

func ImportWithIDAndOrgV5(_ context.Context, data *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	return ImportWithIDAndAttributesV5(ImportOrgAttribute)(context.Background(), data, nil)
}

func ImportWithIDAndOptionalSecretStringV5(secretKey string) schema.StateContextFunc {
	return func(ctx context.Context, data *schema.ResourceData, i interface{}) ([]*schema.ResourceData, error) {
		return ImportWithIDAndAttributesV5(ImportAttribute{Key: secretKey, ValueFromString: ConvertNonEmpty, Optional: true})(context.Background(), data, nil)
	}
}

func ImportWithIDAndOrgAndOptionalSecretStringV5(secretKey string) schema.StateContextFunc {
	return func(ctx context.Context, data *schema.ResourceData, i interface{}) ([]*schema.ResourceData, error) {
		return ImportWithIDAndAttributesV5(ImportOrgAttribute, ImportAttribute{Key: secretKey, ValueFromString: ConvertNonEmpty, Optional: true})(context.Background(), data, nil)
	}
}

func ImportWithOptionalIDV5(idVar string) schema.StateContextFunc {
	return func(ctx context.Context, data *schema.ResourceData, i interface{}) ([]*schema.ResourceData, error) {
		return ImportWithAttributesV5(ImportAttribute{Key: idVar, ValueFromString: ConvertOptionalID, Optional: true})(context.Background(), data, nil)
	}
}

type ImportAttribute struct {
	Key             string
	ValueFromString func(string) (interface{}, error)
	Optional        bool
}

func ImportWithIDAndAttributesV5(attrs ...ImportAttribute) schema.StateContextFunc {
	return ImportWithAttributesV5(append([]ImportAttribute{{Key: "id", Optional: false, ValueFromString: ConvertID}}, attrs...)...)
}

func ImportWithAttributesV5(attrs ...ImportAttribute) schema.StateContextFunc {
	return func(ctx context.Context, data *schema.ResourceData, i interface{}) (ret []*schema.ResourceData, err error) {
		id := data.Id()
		var (
			optionalKeys []string
			requiredKeys []string
		)
		for _, attr := range attrs {
			if attr.Optional {
				optionalKeys = append(optionalKeys, attr.Key)
			} else {
				requiredKeys = append(requiredKeys, attr.Key)
			}
		}
		defer func() {
			if err != nil {
				expectFormat := fmt.Sprintf("<")
				if len(requiredKeys) > 0 {
					expectFormat += fmt.Sprintf(":%s", strings.Join(requiredKeys, ":"))
				}
				if len(optionalKeys) > 0 {
					expectFormat += fmt.Sprintf("[:%s]", strings.Join(optionalKeys, "][:"))
				}
				expectFormat = strings.Replace(expectFormat, "<:", "<", 1)
				expectFormat += ">"
				err = fmt.Errorf("failed to import id %s by format %s: %w", id, expectFormat, err)
			}
		}()
		parts := strings.SplitN(id, ":", len(attrs))
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
				continue
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

func ConvertOptionalID(importValue string) (interface{}, error) {
	if len(importValue) == 0 {
		return "imported", nil
	}
	return ConvertID(importValue)
}

func ConvertNonEmpty(importValue string) (interface{}, error) {
	if len(importValue) == 0 {
		return nil, errors.New("value must not be empty")
	}
	return importValue, nil
}
