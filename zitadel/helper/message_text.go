package helper

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	resourceschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// SMSOTPUnsupportedTextAttrs are the MessageCustomText attributes that the SMS
// OTP message-text API does not persist. Its Set request only carries
// `language` and `text`, so these fields are silently dropped on write and
// returned empty on read. They are shared with the e-mail message-text schema,
// which legitimately uses them, so they stay in the SMS OTP schema (deprecated)
// to keep existing configurations parsing.
var SMSOTPUnsupportedTextAttrs = []string{"title", "pre_header", "subject", "greeting", "button_text", "footer_text"}

const smsOTPUnsupportedDeprecation = "Not supported by the SMS OTP message: the ZITADEL API only stores `text` for this message type, so this attribute has no effect. It will be removed in a future major version."

// DeprecateSMSOTPTextAttrs marks the unsupported MessageCustomText attributes as
// deprecated in place, keeping them optional so existing configurations continue
// to parse without error.
func DeprecateSMSOTPTextAttrs(s *resourceschema.Schema) {
	for _, name := range SMSOTPUnsupportedTextAttrs {
		if _, ok := s.Attributes[name]; !ok {
			continue
		}
		s.Attributes[name] = resourceschema.StringAttribute{
			Optional:           true,
			DeprecationMessage: smsOTPUnsupportedDeprecation,
		}
	}
}

// PreserveSMSOTPTextAttrs copies the unsupported attribute values from src into
// dst. A Read otherwise overwrites the user's configured values with the empty
// values the API returns for these fields, producing a perpetual diff.
func PreserveSMSOTPTextAttrs(ctx context.Context, src types.Object, dst *types.Object) diag.Diagnostics {
	srcAttrs := src.Attributes()
	dstAttrs := dst.Attributes()
	for _, name := range SMSOTPUnsupportedTextAttrs {
		if v, ok := srcAttrs[name]; ok {
			dstAttrs[name] = v
		}
	}

	newObj, d := types.ObjectValue(dst.AttributeTypes(ctx), dstAttrs)
	if d.HasError() {
		return d
	}

	*dst = newObj
	return nil
}
