// Package machine_user_client_secret implements the
// zitadel_machine_user_client_secret ephemeral resource.
//
// It generates a new client secret for a machine (service) user and returns
// the client_id/client_secret pair for the duration of a single apply only,
// never writing the secret to Terraform state (issue #413). Evaluating it
// rotates the secret, so gate it behind count/for_each for on-demand rotation.
package machine_user_client_secret

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
	UserID       types.String `tfsdk:"user_id"`
	OrgID        types.String `tfsdk:"org_id"`
	ClientID     types.String `tfsdk:"client_id"`
	ClientSecret types.String `tfsdk:"client_secret"`
}

func (r *resourceImpl) Metadata(_ context.Context, req ephemeral.MetadataRequest, resp *ephemeral.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_machine_user_client_secret"
}

func (r *resourceImpl) Schema(_ context.Context, _ ephemeral.SchemaRequest, resp *ephemeral.SchemaResponse) {
	resp.Schema = ephschema.Schema{
		MarkdownDescription: "Generates and returns a new client secret for a machine (service) user without persisting " +
			"it to Terraform state. Evaluating this ephemeral resource rotates the secret.",
		Attributes: map[string]ephschema.Attribute{
			"user_id": ephschema.StringAttribute{
				Required:    true,
				Description: "ID of the machine user whose client secret should be generated.",
			},
			"org_id": ephschema.StringAttribute{
				Optional:    true,
				Description: "ID of the organization that owns the machine user. Defaults to the organization of the authenticated user.",
			},
			"client_id": ephschema.StringAttribute{
				Computed:    true,
				Description: "The client ID associated with the generated secret.",
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

	zResp, err := client.GenerateMachineSecret(ctx, &management.GenerateMachineSecretRequest{
		UserId: data.UserID.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("failed to generate machine secret", err.Error())
		return
	}

	data.ClientID = types.StringValue(zResp.GetClientId())
	data.ClientSecret = types.StringValue(zResp.GetClientSecret())
	resp.Diagnostics.Append(resp.Result.Set(ctx, &data)...)
}
