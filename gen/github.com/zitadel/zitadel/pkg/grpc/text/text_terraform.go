package text

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	resourceschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	textpb "github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/text"
	"google.golang.org/protobuf/encoding/protojson"
	proto "google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

const (
	idAttr       = "id"
	orgIDAttr    = "org_id"
	languageAttr = "language"
)

var (
	loginSchema           resourceschema.Schema
	loginAttrTypes        map[string]attr.Type
	messageSchema         resourceschema.Schema
	messageAttrTypes      map[string]attr.Type
	attributeDescriptions = map[string]string{
		"success_login_text.auto_redirect_description": "Text to describe that auto-redirect should happen after successful login",
		"success_login_text.redirected_description":    "Text to describe that the window can be closed after redirect",
	}
	skipProtoRootKeys = map[string]struct{}{
		idAttr:       {},
		orgIDAttr:    {},
		languageAttr: {},
	}
)

func init() {
	loginSchemaAttrs, loginTypes := buildSchemaForMessage((&textpb.LoginCustomText{}).ProtoReflect().Descriptor(), "")
	augmentRootAttributes(loginSchemaAttrs, loginTypes)
	loginSchema = resourceschema.Schema{Attributes: loginSchemaAttrs}
	loginAttrTypes = loginTypes

	messageSchemaAttrs, messageTypes := buildSchemaForMessage((&textpb.MessageCustomText{}).ProtoReflect().Descriptor(), "")
	augmentRootAttributes(messageSchemaAttrs, messageTypes)
	messageSchema = resourceschema.Schema{Attributes: messageSchemaAttrs}
	messageAttrTypes = messageTypes
}

// GenSchemaLoginCustomText returns the schema for LoginCustomText.
func GenSchemaLoginCustomText(_ context.Context) (resourceschema.Schema, diag.Diagnostics) {
	return cloneSchema(loginSchema), nil
}

// GenSchemaMessageCustomText returns the schema for MessageCustomText.
func GenSchemaMessageCustomText(_ context.Context) (resourceschema.Schema, diag.Diagnostics) {
	return cloneSchema(messageSchema), nil
}

// CopyLoginCustomTextFromTerraform populates the proto object from the Terraform value.
func CopyLoginCustomTextFromTerraform(ctx context.Context, tf types.Object, obj *textpb.LoginCustomText) diag.Diagnostics {
	return copyFromTerraform(ctx, tf, obj, loginAttrTypes)
}

// CopyLoginCustomTextToTerraform writes the proto object into the Terraform value.
func CopyLoginCustomTextToTerraform(ctx context.Context, obj *textpb.LoginCustomText, tf *types.Object) diag.Diagnostics {
	return copyToTerraform(ctx, obj, tf, loginAttrTypes)
}

// CopyMessageCustomTextFromTerraform populates the proto object from the Terraform value.
func CopyMessageCustomTextFromTerraform(ctx context.Context, tf types.Object, obj *textpb.MessageCustomText) diag.Diagnostics {
	return copyFromTerraform(ctx, tf, obj, messageAttrTypes)
}

// CopyMessageCustomTextToTerraform writes the proto object into the Terraform value.
func CopyMessageCustomTextToTerraform(ctx context.Context, obj *textpb.MessageCustomText, tf *types.Object) diag.Diagnostics {
	return copyToTerraform(ctx, obj, tf, messageAttrTypes)
}

func buildSchemaForMessage(desc protoreflect.MessageDescriptor, prefix string) (map[string]resourceschema.Attribute, map[string]attr.Type) {
	attrs := make(map[string]resourceschema.Attribute, desc.Fields().Len())
	attrTypes := make(map[string]attr.Type, desc.Fields().Len())

	for i := 0; i < desc.Fields().Len(); i++ {
		field := desc.Fields().Get(i)
		name := string(field.Name())
		path := name
		if prefix != "" {
			path = prefix + "." + name
		}
		if name == "details" {
			continue
		}

		switch field.Kind() {
		case protoreflect.StringKind:
			attrs[name] = resourceschema.StringAttribute{
				Optional:    true,
				Description: attributeDescriptions[path],
			}
			attrTypes[name] = types.StringType
		case protoreflect.MessageKind:
			childAttrs, childTypes := buildSchemaForMessage(field.Message(), path)
			attrs[name] = resourceschema.SingleNestedAttribute{
				Optional:   true,
				Attributes: childAttrs,
			}
			attrTypes[name] = types.ObjectType{AttrTypes: childTypes}
		case protoreflect.BoolKind:
			// Skip internal flags such as is_default to match the previous generated schema.
			continue
		default:
			continue
		}
	}

	return attrs, attrTypes
}

func augmentRootAttributes(attrs map[string]resourceschema.Attribute, attrTypes map[string]attr.Type) {
	attrs[idAttr] = resourceschema.StringAttribute{Computed: true}
	attrs[orgIDAttr] = resourceschema.StringAttribute{Required: true}
	attrs[languageAttr] = resourceschema.StringAttribute{Required: true}

	attrTypes[idAttr] = types.StringType
	attrTypes[orgIDAttr] = types.StringType
	attrTypes[languageAttr] = types.StringType
}

func cloneSchema(schema resourceschema.Schema) resourceschema.Schema {
	schema.Attributes = cloneAttributes(schema.Attributes)
	return schema
}

func cloneAttributes(attrs map[string]resourceschema.Attribute) map[string]resourceschema.Attribute {
	if attrs == nil {
		return nil
	}

	cloned := make(map[string]resourceschema.Attribute, len(attrs))
	for k, v := range attrs {
		cloned[k] = cloneAttribute(v)
	}
	return cloned
}

func cloneAttribute(attr resourceschema.Attribute) resourceschema.Attribute {
	switch a := attr.(type) {
	case resourceschema.SingleNestedAttribute:
		c := a
		c.Attributes = cloneAttributes(a.Attributes)
		return c
	default:
		return attr
	}
}

func copyFromTerraform(ctx context.Context, tf types.Object, obj proto.Message, attrTypes map[string]attr.Type) diag.Diagnostics {
	var diags diag.Diagnostics
	if obj == nil {
		return diags
	}

	data := objectToMap(ctx, tf, attrTypes)
	for key := range skipProtoRootKeys {
		delete(data, key)
	}

	if len(data) == 0 {
		return diags
	}

	raw, err := json.Marshal(data)
	if err != nil {
		diags.AddError("failed to encode text object", err.Error())
		return diags
	}

	if err := (protojson.UnmarshalOptions{DiscardUnknown: true}).Unmarshal(raw, obj); err != nil {
		diags.AddError("failed to decode text object", err.Error())
	}

	return diags
}

func copyToTerraform(ctx context.Context, obj proto.Message, tf *types.Object, attrTypes map[string]attr.Type) diag.Diagnostics {
	var diags diag.Diagnostics
	if tf == nil {
		return diags
	}

	targetAttrTypes := attrTypes
	if existing := tf.AttributeTypes(ctx); len(existing) > 0 {
		targetAttrTypes = existing
	}

	current := objectToMap(ctx, *tf, targetAttrTypes)
	if obj != nil {
		raw, err := (protojson.MarshalOptions{EmitUnpopulated: true, UseProtoNames: true}).Marshal(obj)
		if err != nil {
			diags.AddError("failed to marshal text object", err.Error())
			return diags
		}

		var fromProto map[string]any
		if err := json.Unmarshal(raw, &fromProto); err != nil {
			diags.AddError("failed to decode marshalled text object", err.Error())
			return diags
		}

		for k, v := range fromProto {
			current[k] = v
		}
	}

	values, valueDiags := mapToAttrValues(ctx, current, targetAttrTypes)
	diags.Append(valueDiags...)
	if diags.HasError() {
		return diags
	}

	newObj, objDiags := types.ObjectValue(targetAttrTypes, values)
	diags.Append(objDiags...)
	if diags.HasError() {
		return diags
	}

	*tf = newObj
	return diags
}

func objectToMap(ctx context.Context, obj types.Object, attrTypes map[string]attr.Type) map[string]any {
	if attrTypes == nil {
		return map[string]any{}
	}

	result := make(map[string]any, len(attrTypes))
	attrs := obj.Attributes()

	for name, typ := range attrTypes {
		val, ok := attrs[name]
		if !ok {
			continue
		}

		switch t := typ.(type) {
		case basetypes.StringType:
			if str, ok := val.(types.String); ok && !str.IsNull() && !str.IsUnknown() {
				result[name] = str.ValueString()
			}
		case basetypes.ObjectType:
			objVal, ok := val.(types.Object)
			if !ok || objVal.IsNull() || objVal.IsUnknown() {
				continue
			}
			result[name] = objectToMap(ctx, objVal, t.AttrTypes)
		}
	}

	return result
}

func mapToAttrValues(ctx context.Context, data map[string]any, attrTypes map[string]attr.Type) (map[string]attr.Value, diag.Diagnostics) {
	values := make(map[string]attr.Value, len(attrTypes))
	var diags diag.Diagnostics

	for name, typ := range attrTypes {
		switch t := typ.(type) {
		case basetypes.StringType:
			if raw, ok := data[name]; ok {
				if raw == nil {
					values[name] = types.StringNull()
					break
				}
				if s, ok := raw.(string); ok {
					values[name] = types.StringValue(s)
					break
				}
			}
			values[name] = types.StringNull()
		case basetypes.ObjectType:
			raw, ok := data[name]
			if !ok || raw == nil {
				values[name] = types.ObjectNull(t.AttrTypes)
				break
			}

			rawMap, ok := raw.(map[string]any)
			if !ok {
				values[name] = types.ObjectNull(t.AttrTypes)
				break
			}

			nestedValues, nestedDiags := mapToAttrValues(ctx, rawMap, t.AttrTypes)
			diags.Append(nestedDiags...)
			if diags.HasError() {
				break
			}

			nestedObj, objDiags := types.ObjectValue(t.AttrTypes, nestedValues)
			diags.Append(objDiags...)
			values[name] = nestedObj
		default:
			values[name] = types.StringNull()
		}
	}

	return values, diags
}
