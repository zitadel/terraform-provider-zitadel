package login_texts

import (
	"context"
	"fmt"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	providerschema "github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/management"
	textpb "github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/text"
	"google.golang.org/protobuf/encoding/protojson"

	generatedtext "github.com/zitadel/terraform-provider-zitadel/v2/gen/github.com/zitadel/zitadel/pkg/grpc/text"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
)

const (
	LanguageVar = "language"
)

var (
	_ resource.Resource = &loginTextsResource{}
)

func New() resource.Resource {
	return &loginTextsResource{}
}

type loginTextsResource struct {
	clientInfo *helper.ClientInfo
}

func (r *loginTextsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_login_texts"
}

func (r *loginTextsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	// Get the generated schema - this returns a provider schema, so we need to convert it
	providerSchema, d := generatedtext.GenSchemaLoginCustomText(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Convert provider schema to resource schema
	resourceAttrs := make(map[string]schema.Attribute)
	for name, attr := range providerSchema.Attributes {
		resourceAttrs[name] = convertProviderAttrToResourceAttr(attr)
	}

	resp.Schema = schema.Schema{
		Attributes:          resourceAttrs,
		Description:         providerSchema.Description,
		MarkdownDescription: providerSchema.MarkdownDescription,
		DeprecationMessage:  providerSchema.DeprecationMessage,
	}
}

func (r *loginTextsResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	clientInfo, ok := req.ProviderData.(*helper.ClientInfo)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Configure Provider Data Type",
			fmt.Sprintf("Expected *helper.ClientInfo, got: %T", req.ProviderData),
		)
		return
	}
	r.clientInfo = clientInfo
}

func (r *loginTextsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	orgID, language := getPlanAttrs(ctx, req.Plan, resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	var plan types.Object
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	obj := textpb.LoginCustomText{}
	resp.Diagnostics.Append(generatedtext.CopyLoginCustomTextFromTerraform(ctx, plan, &obj)...)
	if resp.Diagnostics.HasError() {
		return
	}

	jsonpb := &runtime.JSONPb{
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	}
	data, err := jsonpb.Marshal(&obj)
	if err != nil {
		resp.Diagnostics.AddError("failed to marshal protobuf object", err.Error())
		return
	}
	zReq := &management.SetCustomLoginTextsRequest{}
	if err := jsonpb.Unmarshal(data, zReq); err != nil {
		resp.Diagnostics.AddError("failed to unmarshal into ZITADEL request", err.Error())
		return
	}
	zReq.Language = language

	client, err := helper.GetManagementClient(ctx, r.clientInfo)
	if err != nil {
		resp.Diagnostics.AddError("failed to get client", err.Error())
		return
	}

	_, err = client.SetCustomLoginText(helper.CtxSetOrgID(ctx, orgID), zReq)
	if err != nil {
		resp.Diagnostics.AddError("failed to create login texts", err.Error())
		return
	}

	setID(&plan, orgID, language)
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *loginTextsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state types.Object
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	orgID, language := getID(ctx, state)

	client, err := helper.GetManagementClient(ctx, r.clientInfo)
	if err != nil {
		resp.Diagnostics.AddError("failed to get client", err.Error())
		return
	}

	zResp, err := client.GetCustomLoginTexts(helper.CtxSetOrgID(ctx, orgID), &management.GetCustomLoginTextsRequest{Language: language})
	if err != nil {
		if isResourceNotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("failed to read login texts", err.Error())
		return
	}

	if zResp.GetCustomText().GetIsDefault() {
		return
	}

	resp.Diagnostics.Append(generatedtext.CopyLoginCustomTextToTerraform(ctx, zResp.GetCustomText(), &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	setID(&state, orgID, language)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *loginTextsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	orgID, language := getPlanAttrs(ctx, req.Plan, resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	var plan types.Object
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	obj := textpb.LoginCustomText{}
	resp.Diagnostics.Append(generatedtext.CopyLoginCustomTextFromTerraform(ctx, plan, &obj)...)
	if resp.Diagnostics.HasError() {
		return
	}

	jsonpb := &runtime.JSONPb{
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	}
	data, err := jsonpb.Marshal(&obj)
	if err != nil {
		resp.Diagnostics.AddError("failed to marshal protobuf object", err.Error())
		return
	}
	zReq := &management.SetCustomLoginTextsRequest{}
	if err := jsonpb.Unmarshal(data, zReq); err != nil {
		resp.Diagnostics.AddError("failed to unmarshal into ZITADEL request", err.Error())
		return
	}
	zReq.Language = language

	client, err := helper.GetManagementClient(ctx, r.clientInfo)
	if err != nil {
		resp.Diagnostics.AddError("failed to get client", err.Error())
		return
	}

	_, err = client.SetCustomLoginText(helper.CtxSetOrgID(ctx, orgID), zReq)
	if err != nil {
		resp.Diagnostics.AddError("failed to update login texts", err.Error())
		return
	}

	setID(&plan, orgID, language)
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *loginTextsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state types.Object
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	orgID, language := getStateAttrsFromObject(ctx, state, resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := helper.GetManagementClient(ctx, r.clientInfo)
	if err != nil {
		resp.Diagnostics.AddError("failed to get client", err.Error())
		return
	}

	_, err = client.ResetCustomLoginTextToDefault(helper.CtxSetOrgID(ctx, orgID), &management.ResetCustomLoginTextsToDefaultRequest{Language: language})
	if err != nil {
		if isResourceNotFoundError(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("failed to delete login texts", err.Error())
		return
	}
}

func setID(obj *types.Object, orgID string, language string) {
	if obj.IsNull() || obj.IsUnknown() {
		return
	}
	attrs := obj.Attributes()
	if attrs == nil {
		return
	}
	attrs["id"] = types.StringValue(orgID + "_" + language)
	attrs[helper.OrgIDVar] = types.StringValue(orgID)
	attrs[LanguageVar] = types.StringValue(language)
}

func getID(ctx context.Context, obj types.Object) (string, string) {
	if obj.IsNull() || obj.IsUnknown() {
		return "", ""
	}
	id := helper.GetStringFromAttr(ctx, obj.Attributes(), "id")
	parts := strings.Split(id, "_")
	if len(parts) == 2 {
		return parts[0], parts[1]
	}
	return helper.GetStringFromAttr(ctx, obj.Attributes(), helper.OrgIDVar), helper.GetStringFromAttr(ctx, obj.Attributes(), LanguageVar)
}

func getPlanAttrs(ctx context.Context, plan tfsdk.Plan, diags diag.Diagnostics) (string, string) {
	var orgID types.String
	diags.Append(plan.GetAttribute(ctx, path.Root(helper.OrgIDVar), &orgID)...)
	if diags.HasError() || orgID.IsNull() || orgID.IsUnknown() {
		return "", ""
	}
	var language types.String
	diags.Append(plan.GetAttribute(ctx, path.Root(LanguageVar), &language)...)
	if diags.HasError() || language.IsNull() || language.IsUnknown() {
		return "", ""
	}
	return orgID.ValueString(), language.ValueString()
}

func getStateAttrsFromObject(ctx context.Context, obj types.Object, diags diag.Diagnostics) (string, string) {
	if obj.IsNull() || obj.IsUnknown() {
		return "", ""
	}

	attrs := obj.Attributes()
	if attrs == nil {
		return "", ""
	}

	var orgID string
	if orgIDAttr, exists := attrs[helper.OrgIDVar]; exists {
		if orgIDStr, ok := orgIDAttr.(types.String); ok && !orgIDStr.IsNull() && !orgIDStr.IsUnknown() {
			orgID = orgIDStr.ValueString()
		}
	}

	var language string
	if langAttr, exists := attrs[LanguageVar]; exists {
		if langStr, ok := langAttr.(types.String); ok && !langStr.IsNull() && !langStr.IsUnknown() {
			language = langStr.ValueString()
		}
	}

	return orgID, language
}

// Helper function to convert provider schema attributes to resource schema attributes
func convertProviderAttrToResourceAttr(attr providerschema.Attribute) schema.Attribute {
	switch v := attr.(type) {
	case providerschema.StringAttribute:
		return schema.StringAttribute{
			Description:         v.Description,
			MarkdownDescription: v.MarkdownDescription,
			Required:            v.Required,
			Optional:            v.Optional,
			Computed:            false,
			Sensitive:           v.Sensitive,
			DeprecationMessage:  v.DeprecationMessage,
		}
	case providerschema.BoolAttribute:
		return schema.BoolAttribute{
			Description:         v.Description,
			MarkdownDescription: v.MarkdownDescription,
			Required:            v.Required,
			Optional:            v.Optional,
			Computed:            false,
			Sensitive:           v.Sensitive,
			DeprecationMessage:  v.DeprecationMessage,
		}
	case providerschema.Int64Attribute:
		return schema.Int64Attribute{
			Description:         v.Description,
			MarkdownDescription: v.MarkdownDescription,
			Required:            v.Required,
			Optional:            v.Optional,
			Computed:            false,
			Sensitive:           v.Sensitive,
			DeprecationMessage:  v.DeprecationMessage,
		}
	case providerschema.SingleNestedAttribute:
		nestedAttrs := make(map[string]schema.Attribute)
		for name, nestedAttr := range v.Attributes {
			nestedAttrs[name] = convertProviderAttrToResourceAttr(nestedAttr)
		}
		return schema.SingleNestedAttribute{
			Description:         v.Description,
			MarkdownDescription: v.MarkdownDescription,
			Required:            v.Required,
			Optional:            v.Optional,
			Computed:            false,
			Sensitive:           v.Sensitive,
			DeprecationMessage:  v.DeprecationMessage,
			Attributes:          nestedAttrs,
		}
	case providerschema.ListNestedAttribute:
		nestedAttrs := make(map[string]schema.Attribute)
		for name, nestedAttr := range v.NestedObject.Attributes {
			nestedAttrs[name] = convertProviderAttrToResourceAttr(nestedAttr)
		}
		return schema.ListNestedAttribute{
			Description:         v.Description,
			MarkdownDescription: v.MarkdownDescription,
			Required:            v.Required,
			Optional:            v.Optional,
			Computed:            false,
			Sensitive:           v.Sensitive,
			DeprecationMessage:  v.DeprecationMessage,
			NestedObject: schema.NestedAttributeObject{
				Attributes: nestedAttrs,
			},
		}
	default:
		return schema.StringAttribute{
			Optional: true,
		}
	}
}

// Helper function to check if an error indicates a resource was not found
func isResourceNotFoundError(err error) bool {
	if err == nil {
		return false
	}

	errStr := strings.ToLower(err.Error())
	return strings.Contains(errStr, "not found") ||
		strings.Contains(errStr, "not_found") ||
		strings.Contains(errStr, "does not exist")
}