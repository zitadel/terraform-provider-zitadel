package org_idp_github

import (
	"context"

	"google.golang.org/grpc/codes"

	"google.golang.org/grpc/status"

	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/management"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/helper"
)

var (
	_ resource.ResourceWithConfigure = (*githubOnOrgResource)(nil)
)

func NewResource() resource.Resource {
	return &githubOnOrgResource{}
}

type githubOnOrgResource struct {
	clientInfo *helper.ClientInfo
}

func (r *githubOnOrgResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_org_idp_github"
}

func (r *githubOnOrgResource) GetSchema(context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Description: "Resource representing a GitHub IDP on the organization.",
		Attributes:  schemaAttributes(true),
	}, nil
}

func (r *githubOnOrgResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.clientInfo = req.ProviderData.(*helper.ClientInfo)
}

func (r *githubOnOrgResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	state := new(tfModelSensitive)
	resp.Diagnostics.Append(req.Plan.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	client, err := helper.GetManagementClient(r.clientInfo, state.OrgID)
	if err != nil {
		resp.Diagnostics.AddError("getting management client failed", err.Error())
		return
	}
	createResp, err := client.AddGitHubProvider(ctx, state.toPbAddGithubProviderRequest())
	if err != nil {
		resp.Diagnostics.AddError("failed to create idp", err.Error())
		return
	}
	state.ID = createResp.GetId()
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *githubOnOrgResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	oldState := new(tfModelSensitive)
	resp.Diagnostics.Append(req.State.Get(ctx, oldState)...)
	if resp.Diagnostics.HasError() {
		return
	}
	client, err := helper.GetManagementClient(r.clientInfo, oldState.OrgID)
	if err != nil {
		resp.Diagnostics.AddError("getting management client failed", err.Error())
		return
	}
	readResp, err := client.GetProviderByID(ctx, &management.GetProviderByIDRequest{Id: oldState.ID})
	if status.Code(err) == codes.NotFound {
		resp.State.RemoveResource(ctx)
		return
	}
	newState := fromPbGetProviderByIDResponse(readResp, oldState.OrgID).withClientSecret(oldState.ClientSecret.ValueString())
	resp.Diagnostics.Append(resp.State.Set(ctx, newState)...)
}

func (r *githubOnOrgResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	state := new(tfModelSensitive)
	resp.Diagnostics.Append(req.Plan.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	client, err := helper.GetManagementClient(r.clientInfo, state.OrgID)
	if err != nil {
		resp.Diagnostics.AddError("getting management client failed", err.Error())
		return
	}
	_, err = client.UpdateGitHubProvider(ctx, state.toPbUpdateGithubProviderRequest())
	if err != nil {
		resp.Diagnostics.AddError("failed to update idp", err.Error())
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *githubOnOrgResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	state := new(tfModel)
	resp.Diagnostics.Append(req.State.Get(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	client, err := helper.GetManagementClient(r.clientInfo, state.OrgID)
	if err != nil {
		resp.Diagnostics.AddError("getting management client failed", err.Error())
		return
	}
	_, err = client.DeleteProvider(ctx, &management.DeleteProviderRequest{Id: state.ID})
	if status.Code(err) == codes.NotFound {
		err = nil
	}
	if err != nil {
		resp.Diagnostics.AddError("deleting idp failed", err.Error())
	}
}
