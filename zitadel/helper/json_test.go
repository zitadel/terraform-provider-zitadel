package helper

import (
	"testing"
)

func TestNonEmptyJSONObject(t *testing.T) {
	tests := []struct {
		name    string
		value   interface{}
		wantErr bool
	}{
		{
			name:    "object with keys is accepted",
			value:   `{"loginname":{"title":"Welcome"}}`,
			wantErr: false,
		},
		{
			name:    "empty object is rejected",
			value:   `{}`,
			wantErr: true,
		},
		{
			name:    "null is rejected",
			value:   `null`,
			wantErr: true,
		},
		{
			name:    "array is rejected",
			value:   `[{"title":"Welcome"}]`,
			wantErr: true,
		},
		{
			name:    "scalar is rejected",
			value:   `"Welcome"`,
			wantErr: true,
		},
		{
			name:    "invalid json is rejected",
			value:   `{"title":`,
			wantErr: true,
		},
		{
			name:    "non-string value is rejected",
			value:   42,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			diags := NonEmptyJSONObject("translations")(tt.value, nil)
			if diags.HasError() != tt.wantErr {
				t.Errorf("NonEmptyJSONObject() error = %v, wantErr %v", diags, tt.wantErr)
			}
		})
	}
}
