package org_idp_github

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"google.golang.org/grpc/codes"

	"google.golang.org/grpc/status"

	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/management"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/helper"
)

var (
	_ datasource.DataSourceWithConfigure = (*githubOnOrgDataSource)(nil)
)

func NewDataSource() datasource.DataSource {
	return &githubOnOrgDataSource{}
}

type githubOnOrgDataSource struct {
	clientInfo *helper.ClientInfo
}

func (r *githubOnOrgDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_org_idp_github"
}

func (r *githubOnOrgDataSource) GetSchema(context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Description: "Datasource representing a GitHub IDP on the organization.",
		Attributes:  schemaAttributes(false),
	}, nil
}

func (r *githubOnOrgDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.clientInfo = req.ProviderData.(*helper.ClientInfo)
}

func (r *githubOnOrgDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	cfg := new(tfModel)
	resp.Diagnostics.Append(req.Config.Get(ctx, cfg)...)
	if resp.Diagnostics.HasError() {
		return
	}
	client, err := helper.GetManagementClient(r.clientInfo, cfg.OrgID)
	if err != nil {
		resp.Diagnostics.AddError("getting management client failed", err.Error())
		return
	}
	readResp, err := client.GetProviderByID(ctx, &management.GetProviderByIDRequest{Id: cfg.ID})
	if status.Code(err) == codes.NotFound {
		resp.State.RemoveResource(ctx)
		return
	}
	state := fromPbGetProviderByIDResponse(readResp, cfg.OrgID)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}
