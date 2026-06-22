// Package application_api_client_secret implements the
// zitadel_application_api_client_secret ephemeral resource.
//
// It regenerates the client secret of a zitadel_application_api application and
// returns it for the duration of a single apply only, never writing it to
// Terraform state (issue #413). Evaluating it rotates the secret, so gate it
// behind count/for_each for explicit, on-demand rotation.
package application_api_client_secret

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	ephschema "github.com/hashicorp/terraform-plugin-framework/ephemeral/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/management"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
)

var (
	_ ephemeral.EphemeralResource              = &resourceImpl{}
	_ ephemeral.EphemeralResourceWithConfigure = &resourceImpl{}
)

func New() ephemeral.EphemeralResource {
	return &resourceImpl{}
}

type resourceImpl struct {
	helper.EphemeralSecretBase
}

type model struct {
	ProjectID    types.String `tfsdk:"project_id"`
	AppID        types.String `tfsdk:"app_id"`
	OrgID        types.String `tfsdk:"org_id"`
	ClientSecret types.String `tfsdk:"client_secret"`
}

func (r *resourceImpl) Metadata(_ context.Context, req ephemeral.MetadataRequest, resp *ephemeral.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_application_api_client_secret"
}

func (r *resourceImpl) Schema(_ context.Context, _ ephemeral.SchemaRequest, resp *ephemeral.SchemaResponse) {
	resp.Schema = ephschema.Schema{
		MarkdownDescription: "Regenerates and returns the client secret of a `zitadel_application_api` application " +
			"without persisting it to Terraform state. Evaluating this ephemeral resource rotates the secret.",
		Attributes: map[string]ephschema.Attribute{
			"project_id": ephschema.StringAttribute{
				Required:    true,
				Description: "ID of the project the application belongs to.",
			},
			"app_id": ephschema.StringAttribute{
				Required:    true,
				Description: "ID of the application whose client secret should be regenerated.",
			},
			"org_id": ephschema.StringAttribute{
				Optional:    true,
				Description: "ID of the organization that owns the application. Defaults to the organization of the authenticated user.",
			},
			"client_secret": ephschema.StringAttribute{
				Computed:    true,
				Sensitive:   true,
				Description: "The newly generated client secret. Only available during apply; never stored in state.",
			},
		},
	}
}

func (r *resourceImpl) Open(ctx context.Context, req ephemeral.OpenRequest, resp *ephemeral.OpenResponse) {
	var data model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := helper.GetManagementClient(ctx, r.ClientInfo)
	if err != nil {
		resp.Diagnostics.AddError("failed to get client", err.Error())
		return
	}

	if orgID := data.OrgID.ValueString(); orgID != "" {
		ctx = helper.CtxSetOrgID(ctx, orgID)
	}

	zResp, err := client.RegenerateAPIClientSecret(ctx, &management.RegenerateAPIClientSecretRequest{
		ProjectId: data.ProjectID.ValueString(),
		AppId:     data.AppID.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("failed to regenerate client secret", err.Error())
		return
	}

	data.ClientSecret = types.StringValue(zResp.GetClientSecret())
	resp.Diagnostics.Append(resp.Result.Set(ctx, &data)...)
}
