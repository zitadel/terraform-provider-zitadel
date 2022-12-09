package verify_email_message_text

import (
	"context"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/management"
	"google.golang.org/protobuf/encoding/protojson"

	textpb "github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/text"

	"github.com/zitadel/terraform-provider-zitadel/gen/github.com/zitadel/zitadel/pkg/grpc/text"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/helper"
)

const (
	orgIDVar    = "org_id"
	languageVar = "language"
)

var (
	_ resource.Resource = &verifyEmailMessageTextResource{}
)

func New() resource.Resource {
	return &verifyEmailMessageTextResource{}
}

type verifyEmailMessageTextResource struct {
	clientInfo *helper.ClientInfo
}

func (r *verifyEmailMessageTextResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_verify_email_message_text"
}

func (r *verifyEmailMessageTextResource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return text.GenSchemaMessageCustomText(ctx)
}

func (r *verifyEmailMessageTextResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.clientInfo = req.ProviderData.(*helper.ClientInfo)
}

func (r *verifyEmailMessageTextResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	orgID, language := getPlanAttrs(ctx, req.Plan, resp.Diagnostics)
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
	zReq := &management.SetCustomVerifyEmailMessageTextRequest{}
	if err := jsonpb.Unmarshal(data, zReq); err != nil {
		resp.Diagnostics.AddError("failed to unmarshal", err.Error())
		return
	}
	zReq.Language = language

	client, err := helper.GetManagementClient(r.clientInfo, orgID)
	if err != nil {
		resp.Diagnostics.AddError("failed to get client", err.Error())
		return
	}

	_, err = client.SetCustomVerifyEmailMessageText(ctx, zReq)
	if err != nil {
		resp.Diagnostics.AddError("failed to create", err.Error())
		return
	}

	setID(plan, orgID, language)
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *verifyEmailMessageTextResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state types.Object
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	orgID, language := getID(ctx, state)

	client, err := helper.GetManagementClient(r.clientInfo, orgID)
	if err != nil {
		resp.Diagnostics.AddError("failed to get client", err.Error())
		return
	}

	zResp, err := client.GetCustomVerifyEmailMessageText(ctx, &management.GetCustomVerifyEmailMessageTextRequest{Language: language})
	if err != nil {
		return
	}
	if zResp.CustomText.IsDefault {
		return
	}

	resp.Diagnostics.Append(text.CopyMessageCustomTextToTerraform(ctx, *zResp.CustomText, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	setID(state, orgID, language)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *verifyEmailMessageTextResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	orgID, language := getPlanAttrs(ctx, req.Plan, resp.Diagnostics)
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
	zReq := &management.SetCustomVerifyEmailMessageTextRequest{}
	if err := jsonpb.Unmarshal(data, zReq); err != nil {
		resp.Diagnostics.AddError("failed to unmarshal", err.Error())
		return
	}
	zReq.Language = language

	client, err := helper.GetManagementClient(r.clientInfo, orgID)
	if err != nil {
		resp.Diagnostics.AddError("failed to get client", err.Error())
		return
	}

	_, err = client.SetCustomVerifyEmailMessageText(ctx, zReq)
	if err != nil {
		resp.Diagnostics.AddError("failed to update", err.Error())
		return
	}

	setID(plan, orgID, language)
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *verifyEmailMessageTextResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	orgID, language := getStateAttrs(ctx, req.State, resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := helper.GetManagementClient(r.clientInfo, orgID)
	if err != nil {
		resp.Diagnostics.AddError("failed to get client", err.Error())
		return
	}

	_, err = client.ResetCustomVerifyEmailMessageTextToDefault(ctx, &management.ResetCustomVerifyEmailMessageTextToDefaultRequest{Language: language})
	if err != nil {
		resp.Diagnostics.AddError("failed to delete", err.Error())
		return
	}
}

func setID(obj types.Object, orgID string, language string) {
	attrs := obj.Attributes()
	attrs["id"] = types.StringValue(orgID + "_" + language)
	attrs[orgIDVar] = types.StringValue(orgID)
	attrs[languageVar] = types.StringValue(language)
}

func getID(ctx context.Context, obj types.Object) (string, string) {
	id := helper.GetStringFromAttr(ctx, obj.Attributes(), "id")
	parts := strings.Split(id, "_")
	if len(parts) == 2 {
		return parts[0], parts[1]
	}
	return helper.GetStringFromAttr(ctx, obj.Attributes(), orgIDVar), helper.GetStringFromAttr(ctx, obj.Attributes(), languageVar)
}

func getPlanAttrs(ctx context.Context, plan tfsdk.Plan, diag diag.Diagnostics) (string, string) {
	var orgID string
	diag.Append(plan.GetAttribute(ctx, path.Root(orgIDVar), &orgID)...)
	if diag.HasError() {
		return "", ""
	}
	var language string
	diag.Append(plan.GetAttribute(ctx, path.Root(languageVar), &language)...)
	if diag.HasError() {
		return "", ""
	}

	return orgID, language
}

func getStateAttrs(ctx context.Context, state tfsdk.State, diag diag.Diagnostics) (string, string) {
	var orgID string
	diag.Append(state.GetAttribute(ctx, path.Root(orgIDVar), &orgID)...)
	if diag.HasError() {
		return "", ""
	}
	var language string
	diag.Append(state.GetAttribute(ctx, path.Root(languageVar), &language)...)
	if diag.HasError() {
		return "", ""
	}

	return orgID, language
}
