package text

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	textpb "github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/text"
)

// TestCopyLoginCustomTextToTerraform verifies how nested screen-text blocks are
// written into state. ZITADEL returns every screen-text message on
// GetCustomLoginTexts, even for screens that were never customised, as a
// present message with all fields empty. Such a block must stay null when the
// prior state does not hold it, otherwise a config that omits it diffs against
// state on every plan. A block the prior state does hold, such as one
// configured as {}, must be kept so the opposite diff does not appear either.
func TestCopyLoginCustomTextToTerraform(t *testing.T) {
	const attribute = "password_change_done_text"
	attrTypes := loginAttrTypes[attribute].(types.ObjectType).AttrTypes
	tests := []struct {
		name      string
		prior     attr.Value
		obj       *textpb.LoginCustomText
		wantNull  bool
		wantTitle string
	}{
		{
			name:     "absent message stays null",
			prior:    types.ObjectNull(attrTypes),
			obj:      &textpb.LoginCustomText{},
			wantNull: true,
		},
		{
			name:  "empty message stays null",
			prior: types.ObjectNull(attrTypes),
			obj: &textpb.LoginCustomText{
				PasswordChangeDoneText: &textpb.PasswordChangeDoneScreenText{},
			},
			wantNull: true,
		},
		{
			name: "empty message keeps explicitly empty block",
			prior: types.ObjectValueMust(attrTypes, map[string]attr.Value{
				"title":            types.StringNull(),
				"description":      types.StringNull(),
				"next_button_text": types.StringNull(),
			}),
			obj: &textpb.LoginCustomText{
				PasswordChangeDoneText: &textpb.PasswordChangeDoneScreenText{},
			},
			wantNull: false,
		},
		{
			name:  "populated message is copied",
			prior: types.ObjectNull(attrTypes),
			obj: &textpb.LoginCustomText{
				PasswordChangeDoneText: &textpb.PasswordChangeDoneScreenText{Title: "Password changed"},
			},
			wantTitle: "Password changed",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			state := loginState(t, attribute, tt.prior)
			if diags := CopyLoginCustomTextToTerraform(context.Background(), tt.obj, &state); diags.HasError() {
				t.Fatalf("unexpected diagnostics: %v", diags)
			}
			got := state.Attributes()[attribute].(types.Object)
			if got.IsNull() != tt.wantNull {
				t.Fatalf("%s null = %v, want %v", attribute, got.IsNull(), tt.wantNull)
			}
			if tt.wantNull {
				return
			}
			if title := got.Attributes()["title"].(types.String); title.ValueString() != tt.wantTitle {
				t.Errorf("%s.title = %q, want %q", attribute, title.ValueString(), tt.wantTitle)
			}
		})
	}
}

// loginState returns a login texts object with every attribute null except the
// given one, which is set to value.
func loginState(t *testing.T, name string, value attr.Value) types.Object {
	t.Helper()
	values := make(map[string]attr.Value, len(loginAttrTypes))
	for attrName, typ := range loginAttrTypes {
		switch typ := typ.(type) {
		case types.ObjectType:
			values[attrName] = types.ObjectNull(typ.AttrTypes)
		default:
			values[attrName] = types.StringNull()
		}
	}
	values[name] = value
	return types.ObjectValueMust(loginAttrTypes, values)
}
