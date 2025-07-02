package default_password_change_message_text

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/admin"
	textpb "github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/text"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/zitadel/terraform-provider-zitadel/v2/gen/github.com/zitadel/zitadel/pkg/grpc/text"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
)

const (
	LanguageVar = "language"
)

var (
	_ resource.Resource = &defaultPasswordChangeMessageTextResource{}
)

func New() resource.Resource {
	return &defaultPasswordChangeMessageTextResource{}
}

type defaultPasswordChangeMessageTextResource struct {
	clientInfo *helper.ClientInfo
}

func (r *defaultPasswordChangeMessageTextResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_password_change_message_text"
}

// Fixed Schema method - properly handle the schema conversion
func (r *defaultPasswordChangeMessageTextResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	// Get the generated schema - this likely returns a provider schema
	generatedSchema, d := text.GenSchemaMessageCustomText(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Convert to resource schema by copying the attributes
	resourceAttrs := make(map[string]schema.Attribute)
	for name, attr := range generatedSchema.Attributes {
		// Copy each attribute, converting from provider schema to resource schema
		// This is a simplified conversion - you may need to handle specific attribute types
		resourceAttrs[name] = convertAttribute(attr)
	}
	
	// Remove org_id if it exists
	delete(resourceAttrs, "org_id")

	resp.Schema = schema.Schema{
		Attributes:  resourceAttrs,
		Description: generatedSchema.Description,
		// Copy other schema properties as needed
	}
}

// Helper function to convert provider schema attributes to resource schema attributes
// You'll need to implement this based on the actual attribute types used
func convertAttribute(providerAttr interface{}) schema.Attribute {
	// This is a placeholder - you'll need to implement the actual conversion
	// based on the specific attribute types returned by text.GenSchemaMessageCustomText
	
	// For example, if it's a string attribute:
	// if stringAttr, ok := providerAttr.(provider_schema.StringAttribute); ok {
	//     return schema.StringAttribute{
	//         Description: stringAttr.Description,
	//         Required:    stringAttr.Required,
	//         Optional:    stringAttr.Optional,
	//         // ... other properties
	//     }
	// }
	
	// For now, return a basic string attribute - replace with proper conversion
	return schema.StringAttribute{
		Optional: true,
	}
}

func (r *defaultPasswordChangeMessageTextResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.clientInfo = req.ProviderData.(*helper.ClientInfo)
}

func (r *defaultPasswordChangeMessageTextResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	language := getPlanAttrs(ctx, req.Plan, resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	var plan types.Object
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	obj := textpb.MessageCustomText{}
	resp.Diagnostics.Append(text.CopyMessageCustomTextFromTerraform(ctx, plan, &obj)...)
	if resp.Diagnostics.HasError() {
		return
	}

	jsonpb := &runtime.JSONPb{
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	}
	data, err := jsonpb.Marshal(obj)
	if err != nil {
		resp.Diagnostics.AddError("failed to marshal", err.Error())
		return
	}
	zReq := &admin.SetDefaultPasswordChangeMessageTextRequest{}
	if err := jsonpb.Unmarshal(data, zReq); err != nil {
		resp.Diagnostics.AddError("failed to unmarshal", err.Error())
		return
	}
	zReq.Language = language

	client, err := helper.GetAdminClient(ctx, r.clientInfo)
	if err != nil {
		resp.Diagnostics.AddError("failed to get client", err.Error())
		return
	}

	_, err = client.SetDefaultPasswordChangeMessageText(ctx, zReq)
	if err != nil {
		resp.Diagnostics.AddError("failed to create", err.Error())
		return
	}

	setID(plan, language)
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *defaultPasswordChangeMessageTextResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state types.Object
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	language := getID(ctx, state)

	client, err := helper.GetAdminClient(ctx, r.clientInfo)
	if err != nil {
		resp.Diagnostics.AddError("failed to get client", err.Error())
		return
	}

	zResp, err := client.GetCustomPasswordChangeMessageText(ctx, &admin.GetCustomPasswordChangeMessageTextRequest{Language: language})
	if err != nil {
		return
	}
	if zResp.CustomText.IsDefault {
		return
	}

	resp.Diagnostics.Append(text.CopyMessageCustomTextToTerraform(ctx, zResp.CustomText, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	setID(state, language)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *defaultPasswordChangeMessageTextResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	language := getPlanAttrs(ctx, req.Plan, resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	var plan types.Object
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	obj := textpb.MessageCustomText{}
	resp.Diagnostics.Append(text.CopyMessageCustomTextFromTerraform(ctx, plan, &obj)...)
	if resp.Diagnostics.HasError() {
		return
	}

	jsonpb := &runtime.JSONPb{
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	}
	data, err := jsonpb.Marshal(obj)
	if err != nil {
		resp.Diagnostics.AddError("failed to marshal", err.Error())
		return
	}
	zReq := &admin.SetDefaultPasswordChangeMessageTextRequest{}
	if err := jsonpb.Unmarshal(data, zReq); err != nil {
		resp.Diagnostics.AddError("failed to unmarshal", err.Error())
		return
	}
	zReq.Language = language

	client, err := helper.GetAdminClient(ctx, r.clientInfo)
	if err != nil {
		resp.Diagnostics.AddError("failed to get client", err.Error())
		return
	}

	_, err = client.SetDefaultPasswordChangeMessageText(ctx, zReq)
	if err != nil {
		resp.Diagnostics.AddError("failed to update", err.Error())
		return
	}

	setID(plan, language)
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *defaultPasswordChangeMessageTextResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	language := getStateAttrs(ctx, req.State, resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := helper.GetAdminClient(ctx, r.clientInfo)
	if err != nil {
		resp.Diagnostics.AddError("failed to get client", err.Error())
		return
	}

	_, err = client.ResetCustomPasswordChangeMessageTextToDefault(ctx, &admin.ResetCustomPasswordChangeMessageTextToDefaultRequest{Language: language})
	if err != nil {
		resp.Diagnostics.AddError("failed to delete", err.Error())
		return
	}
}

func setID(obj types.Object, language string) {
	attrs := obj.Attributes()
	attrs["id"] = types.StringValue(language)
	attrs[LanguageVar] = types.StringValue(language)
}

func getID(ctx context.Context, obj types.Object) string {
	return helper.GetStringFromAttr(ctx, obj.Attributes(), "id")
}

func getPlanAttrs(ctx context.Context, plan tfsdk.Plan, diag diag.Diagnostics) string {
	var language string
	diag.Append(plan.GetAttribute(ctx, path.Root(LanguageVar), &language)...)
	if diag.HasError() {
		return ""
	}
	return language
}

func getStateAttrs(ctx context.Context, state tfsdk.State, diag diag.Diagnostics) string {
	var language string
	diag.Append(state.GetAttribute(ctx, path.Root(LanguageVar), &language)...)
	if diag.HasError() {
		return ""
	}
	return language
}