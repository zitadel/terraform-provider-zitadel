package default_invite_user_message_text

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
	_ resource.Resource = &defaultInviteUserMessageTextResource{}
)

func New() resource.Resource {
	return &defaultInviteUserMessageTextResource{}
}

type defaultInviteUserMessageTextResource struct {
	clientInfo *helper.ClientInfo
}

func (r *defaultInviteUserMessageTextResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_Invite_user_message_text"
}

func (r *defaultInviteUserMessageTextResource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	s, d := text.GenSchemaMessageCustomText(ctx)
	delete(s.Attributes, "org_id")
	return s, d
}

func (r *defaultInviteUserMessageTextResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.clientInfo = req.ProviderData.(*helper.ClientInfo)
}

func (r *defaultInviteUserMessageTextResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
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
	zReq := &admin.SetDefaultInviteUserMessageTextRequest{}
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

	_, err = client.SetDefaultInviteUserMessageText(ctx, zReq)
	if err != nil {
		resp.Diagnostics.AddError("failed to create", err.Error())
		return
	}

	setID(plan, language)
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *defaultInviteUserMessageTextResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
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

	zResp, err := client.GetCustomInviteUserMessageText(ctx, &admin.GetCustomInviteUserMessageTextRequest{Language: language})
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

func (r *defaultInviteUserMessageTextResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
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
	zReq := &admin.SetDefaultInviteUserMessageTextRequest{}
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

	_, err = client.SetDefaultInviteUserMessageText(ctx, zReq)
	if err != nil {
		resp.Diagnostics.AddError("failed to update", err.Error())
		return
	}

	setID(plan, language)
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *defaultInviteUserMessageTextResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	language := getStateAttrs(ctx, req.State, resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := helper.GetAdminClient(ctx, r.clientInfo)
	if err != nil {
		resp.Diagnostics.AddError("failed to get client", err.Error())
		return
	}

	_, err = client.ResetCustomInviteUserMessageTextToDefault(ctx, &admin.ResetCustomInviteUserMessageTextToDefaultRequest{Language: language})
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
