package helper_test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	resourceschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
)

// TestDeprecateSMSOTPTextAttrs verifies that the unsupported attributes gain a
// deprecation message while their generated metadata is preserved, and that
// supported attributes (text) are left untouched.
func TestDeprecateSMSOTPTextAttrs(t *testing.T) {
	s := &resourceschema.Schema{
		Attributes: map[string]resourceschema.Attribute{
			"text":     resourceschema.StringAttribute{Optional: true, Description: "the text"},
			"greeting": resourceschema.StringAttribute{Optional: true, Description: "the greeting"},
		},
	}

	helper.DeprecateSMSOTPTextAttrs(s)

	greeting, ok := s.Attributes["greeting"].(resourceschema.StringAttribute)
	if !ok {
		t.Fatalf("greeting is not a StringAttribute")
	}
	if greeting.DeprecationMessage == "" {
		t.Errorf("greeting should be deprecated")
	}
	if !greeting.Optional || greeting.Description != "the greeting" {
		t.Errorf("greeting metadata not preserved: %+v", greeting)
	}

	text := s.Attributes["text"].(resourceschema.StringAttribute)
	if text.DeprecationMessage != "" {
		t.Errorf("text must not be deprecated")
	}
}

// TestPreserveSMSOTPTextAttrs verifies that the unsupported attribute values are
// copied from the prior state into the destination (so a Read does not overwrite
// the user's configured values), while supported attributes keep the values that
// were written from the server response.
func TestPreserveSMSOTPTextAttrs(t *testing.T) {
	ctx := context.Background()
	attrTypes := map[string]attr.Type{
		"text":     types.StringType,
		"greeting": types.StringType,
		"subject":  types.StringType,
	}

	prior := types.ObjectValueMust(attrTypes, map[string]attr.Value{
		"text":     types.StringValue("configured text"),
		"greeting": types.StringValue("Greeting"),
		"subject":  types.StringValue("Subject"),
	})

	// dst simulates the object after the server copy: text taken from the API,
	// the unsupported fields reset to null.
	dst := types.ObjectValueMust(attrTypes, map[string]attr.Value{
		"text":     types.StringValue("server text"),
		"greeting": types.StringNull(),
		"subject":  types.StringNull(),
	})

	if diags := helper.PreserveSMSOTPTextAttrs(ctx, prior, &dst); diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}

	got := dst.Attributes()
	if v := got["greeting"].(types.String); v.ValueString() != "Greeting" {
		t.Errorf("greeting not preserved, got %q", v.ValueString())
	}
	if v := got["subject"].(types.String); v.ValueString() != "Subject" {
		t.Errorf("subject not preserved, got %q", v.ValueString())
	}
	if v := got["text"].(types.String); v.ValueString() != "server text" {
		t.Errorf("text should keep the server value, got %q", v.ValueString())
	}
}
