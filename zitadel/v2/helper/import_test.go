package helper

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestImportWithAttributes(t *testing.T) {
	validID := "123456789012345678"
	type args struct {
		attrs []importAttribute
		id    string
	}
	type want struct {
		attributes              map[string]interface{}
		expectErrorWithIDFormat string
		expectErrorWithMinParts int
		expectErrorWithMaxParts int
	}
	tests := []struct {
		name string
		args args
		want want
	}{{
		name: `<id> with '123...' works`,
		args: args{
			attrs: []importAttribute{NewImportAttribute("id", ConvertID, false)},
			id:    validID,
		},
		want: want{
			attributes: map[string]interface{}{
				"id": validID,
			},
		},
	}, {
		name: `<id> with '' fails`,
		args: args{
			attrs: []importAttribute{NewImportAttribute("id", ConvertID, false)},
		},
		want: want{
			expectErrorWithIDFormat: "<id>",
			expectErrorWithMinParts: 1,
			expectErrorWithMaxParts: 1,
		},
	}, {
		name: `<id:required_id> with '123...:123...' works`,
		args: args{
			attrs: []importAttribute{
				NewImportAttribute("id", ConvertID, false),
				NewImportAttribute("required_id", ConvertID, false),
			},
			id: concat(validID, validID),
		},
		want: want{
			attributes: map[string]interface{}{
				"id":          validID,
				"required_id": validID,
			},
		},
	}, {
		name: `<id:required_id> with '123...' fails`,
		args: args{
			attrs: []importAttribute{
				NewImportAttribute("id", ConvertID, false),
				NewImportAttribute("required_id", ConvertID, false),
			},
			id: validID,
		},
		want: want{
			expectErrorWithIDFormat: "<id:required_id>",
			expectErrorWithMinParts: 2,
			expectErrorWithMaxParts: 2,
		},
	}, {
		name: `<id:required_id[:optional_id]> with '123...:123...:123...' works`,
		args: args{
			attrs: []importAttribute{
				NewImportAttribute("id", ConvertID, false),
				NewImportAttribute("required_id", ConvertID, false),
				NewImportAttribute("optional_id", ConvertID, true),
			},
			id: concat(validID, validID, validID),
		},
		want: want{
			attributes: map[string]interface{}{
				"id":          validID,
				"required_id": validID,
				"optional_id": validID,
			},
		},
	}, {
		name: `<id:required_id[:optional_id]> with '123...:123...' works`,
		args: args{
			attrs: []importAttribute{
				NewImportAttribute("id", ConvertID, false),
				NewImportAttribute("required_id", ConvertID, false),
				NewImportAttribute("optional_id", ConvertID, true),
			},
			id: concat(validID, validID),
		},
		want: want{
			attributes: map[string]interface{}{
				"id":          validID,
				"required_id": validID,
			},
		},
	}, {
		name: `<id:required_id[:optional_id]> with '123...:123...:123...:123...' fails`,
		args: args{
			attrs: []importAttribute{
				NewImportAttribute("id", ConvertID, false),
				NewImportAttribute("required_id", ConvertID, false),
				NewImportAttribute("optional_id", ConvertID, true),
			},
			id: concat(validID, validID, validID, validID),
		},
		want: want{
			expectErrorWithIDFormat: "<id:required_id[:optional_id]>",
			expectErrorWithMinParts: 2,
			expectErrorWithMaxParts: 3,
		},
	}, {
		name: `<> with '' works`,
		args: args{
			attrs: []importAttribute{emptyIDAttribute},
		},
		want: want{
			attributes: map[string]interface{}{
				"id": "imported",
			},
		},
	}, {
		name: `<> with '123...' fails`,
		args: args{
			attrs: []importAttribute{emptyIDAttribute},
			id:    validID,
		},
		want: want{
			expectErrorWithIDFormat: "<>",
			expectErrorWithMinParts: -1,
			expectErrorWithMaxParts: -1,
		},
	}, {
		name: `<[org_id]> with '123...' works`,
		args: args{
			attrs: []importAttribute{ImportOptionalOrgAttribute},
			id:    validID,
		},
		want: want{
			attributes: map[string]interface{}{
				"id":     validID,
				"org_id": validID,
			},
		},
	}, {
		name: `<[org_id]> with '' works`,
		args: args{
			attrs: []importAttribute{ImportOptionalOrgAttribute},
		},
		want: want{
			attributes: map[string]interface{}{
				"id": "imported",
			},
		},
	}, {
		name: `<[org_id]> with 'invalid id' fails`,
		args: args{
			attrs: []importAttribute{ImportOptionalOrgAttribute},
			id:    "invalid id",
		},
		want: want{
			expectErrorWithIDFormat: "<[org_id]>",
			expectErrorWithMinParts: -1,
			expectErrorWithMaxParts: -1,
		},
	}, {
		name: `<required_id[:optional_id]> with empty id and '123...:123...' works`,
		args: args{
			attrs: []importAttribute{
				emptyIDAttribute,
				NewImportAttribute("required_id", ConvertID, false),
				NewImportAttribute("optional_id", ConvertID, true),
			},
			id: concat(validID, validID),
		},
		want: want{
			attributes: map[string]interface{}{
				"id":          "imported",
				"required_id": validID,
				"optional_id": validID,
			},
		},
	}, {
		name: `<required_id[:optional_id]> with empty id and '123...:123...:123...' fails`,
		args: args{
			attrs: []importAttribute{
				emptyIDAttribute,
				NewImportAttribute("required_id", ConvertID, false),
				NewImportAttribute("optional_id", ConvertID, true),
			},
			id: concat(validID, validID, validID),
		},
		want: want{
			expectErrorWithIDFormat: "<required_id[:optional_id]>",
			expectErrorWithMinParts: 1,
			expectErrorWithMaxParts: 2,
		},
	}, {
		name: `<required_id[:optional_id]> with empty id and '' fails`,
		args: args{
			attrs: []importAttribute{
				emptyIDAttribute,
				NewImportAttribute("required_id", ConvertID, false),
				NewImportAttribute("optional_id", ConvertID, true),
			},
		},
		want: want{
			expectErrorWithIDFormat: "<required_id[:optional_id]>",
			expectErrorWithMinParts: 1,
			expectErrorWithMaxParts: 2,
		},
	}, {
		name: `<required_id:another_required_id[:optional_id]> with empty id and '123...' fails`,
		args: args{
			attrs: []importAttribute{
				emptyIDAttribute,
				NewImportAttribute("required_id", ConvertID, false),
				NewImportAttribute("another_required_id", ConvertID, false),
				NewImportAttribute("optional_id", ConvertID, true),
			},
			id: validID,
		},
		want: want{
			expectErrorWithIDFormat: "<required_id:another_required_id[:optional_id]>",
			expectErrorWithMinParts: 2,
			expectErrorWithMaxParts: 3,
		},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			state := newMockState()
			state.SetId(tt.args.id)
			err := importWithAttributes(state, tt.args.attrs...)
			wantAttributes := tt.want.attributes
			if err != nil {
				if tt.want.expectErrorWithIDFormat == "" {
					t.Fatalf("importWithAttributes() error = %v, want %v", err, wantAttributes)
				}
				expectBetweenError := fmt.Sprintf("between %d and %d", tt.want.expectErrorWithMinParts, tt.want.expectErrorWithMaxParts)
				if (tt.want.expectErrorWithMinParts > -1 || tt.want.expectErrorWithMaxParts > -1) &&
					!strings.Contains(err.Error(), expectBetweenError) {
					t.Errorf(`expected error to contain "%s", got: %v`, expectBetweenError, err)
				}
				if !strings.Contains(err.Error(), tt.want.expectErrorWithIDFormat) {
					t.Errorf("expected error to contain the expected format '%s', got: %v", tt.want.expectErrorWithIDFormat, err)
				}
				return
			}
			if tt.want.expectErrorWithIDFormat != "" {
				t.Fatalf("expected error with format '%s', got state: %v", tt.want.expectErrorWithIDFormat, state)
			}
			if !reflect.DeepEqual(state, mockState(wantAttributes)) {
				t.Errorf("importWithAttributes() = %v, want %v", state, wantAttributes)
			}
		})
	}
}

func newMockState() mockState { return make(map[string]interface{}) }

type mockState map[string]interface{}

// SetId sets the ID of the state.
func (m mockState) SetId(id string) {
	m["id"] = id
}

// Id returns the ID of the state.
func (m mockState) Id() string {
	return m["id"].(string)
}

// Set sets the value of the given attribute.
func (m mockState) Set(key string, value interface{}) error {
	m[key] = value
	return nil
}

func concat(attr ...string) string {
	return strings.Join(attr, ":")
}
