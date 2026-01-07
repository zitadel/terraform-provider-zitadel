package default_login_texts

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
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
	_ resource.Resource = &defaultLoginTextsResource{}
)

func New() resource.Resource {
	return &defaultLoginTextsResource{}
}

type defaultLoginTextsResource struct {
	clientInfo *helper.ClientInfo
}

func (r *defaultLoginTextsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_login_texts"
}

func (r *defaultLoginTextsResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	s, diags := text.GenSchemaLoginCustomText(ctx)
	resp.Diagnostics.Append(diags...)
	delete(s.Attributes, "org_id")
	resp.Schema = s
}

func (r *defaultLoginTextsResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.clientInfo = req.ProviderData.(*helper.ClientInfo)
}

func (r *defaultLoginTextsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	language := getPlanAttrs(ctx, req.Plan, resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	var plan types.Object
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	obj := textpb.LoginCustomText{}
	resp.Diagnostics.Append(text.CopyLoginCustomTextFromTerraform(ctx, plan, &obj)...)
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
		resp.Diagnostics.AddError("failed to marshal login texts", err.Error())
		return
	}
	zReq := &admin.SetCustomLoginTextsRequest{}
	if err := jsonpb.Unmarshal(data, zReq); err != nil {
		resp.Diagnostics.AddError("failed to unmarshal login texts", err.Error())
		return
	}
	zReq.Language = language

	client, err := helper.GetAdminClient(ctx, r.clientInfo)
	if err != nil {
		resp.Diagnostics.AddError("failed to get client", err.Error())
		return
	}

	_, err = client.SetCustomLoginText(ctx, zReq)
	if err != nil {
		resp.Diagnostics.AddError("failed to create login texts", err.Error())
		return
	}

	setID(plan, language)
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *defaultLoginTextsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
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

	zResp, err := client.GetCustomLoginTexts(ctx, &admin.GetCustomLoginTextsRequest{Language: language})
	if err != nil {
		return
	}
	if zResp.CustomText.IsDefault {
		return
	}

	resp.Diagnostics.Append(text.CopyLoginCustomTextToTerraform(ctx, zResp.CustomText, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	setID(state, language)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *defaultLoginTextsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	language := getPlanAttrs(ctx, req.Plan, resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	var plan types.Object
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	obj := textpb.LoginCustomText{}
	resp.Diagnostics.Append(text.CopyLoginCustomTextFromTerraform(ctx, plan, &obj)...)
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
		resp.Diagnostics.AddError("failed to marshal login texts", err.Error())
		return
	}
	zReq := &admin.SetCustomLoginTextsRequest{}
	if err := jsonpb.Unmarshal(data, zReq); err != nil {
		resp.Diagnostics.AddError("failed to unmarshal login texts", err.Error())
		return
	}
	zReq.Language = language

	client, err := helper.GetAdminClient(ctx, r.clientInfo)
	if err != nil {
		resp.Diagnostics.AddError("failed to get client", err.Error())
		return
	}

	_, err = client.SetCustomLoginText(ctx, zReq)
	if err != nil {
		resp.Diagnostics.AddError("failed to update login texts", err.Error())
		return
	}

	setID(plan, language)
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *defaultLoginTextsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	language := getStateAttrs(ctx, req.State, resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := helper.GetAdminClient(ctx, r.clientInfo)
	if err != nil {
		resp.Diagnostics.AddError("failed to get client", err.Error())
		return
	}

	_, err = client.ResetCustomLoginTextToDefault(ctx, &admin.ResetCustomLoginTextsToDefaultRequest{Language: language})
	if err != nil {
		resp.Diagnostics.AddError("failed to delete login texts", err.Error())
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
