// Package application_v2_client_secret implements the
// zitadel_application_v2_client_secret ephemeral resource.
//
// It rotates (regenerates) the client secret of a zitadel_application_v2 OIDC
// or API application and returns the freshly generated secret for the duration
// of a single `terraform apply` only. Unlike the Computed `client_secret`
// attribute on the managed zitadel_application_v2 resource, the value produced
// here is never written to Terraform state, which keeps the secret out of
// remote state backends (issue #413).
//
// Because opening this ephemeral resource calls the server's GenerateClientSecret
// RPC, every apply that evaluates it rotates the secret. Operators are expected
// to gate it behind count/for_each so rotation is explicit rather than
// happening on every run.
package application_v2_client_secret

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	ephschema "github.com/hashicorp/terraform-plugin-framework/ephemeral/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	appv2pb "github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/application/v2"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
)

// Interface assertions: this resource opens (Open) and needs provider data
// (Configure, provided by the embedded EphemeralSecretBase).
var (
	_ ephemeral.EphemeralResource              = &resourceImpl{}
	_ ephemeral.EphemeralResourceWithConfigure = &resourceImpl{}
)

func New() ephemeral.EphemeralResource {
	return &resourceImpl{}
}

type resourceImpl struct {
	// Configure (from the embedded base) populates ClientInfo before Open.
	helper.EphemeralSecretBase
}

// model maps the HCL config and result attributes to Go. application_id and
// project_id identify the v2 application; org_id optionally scopes the call to
// a specific organization. client_secret is the result populated by Open.
type model struct {
	ApplicationID types.String `tfsdk:"application_id"`
	ProjectID     types.String `tfsdk:"project_id"`
	OrgID         types.String `tfsdk:"org_id"`
	ClientSecret  types.String `tfsdk:"client_secret"`
}

func (r *resourceImpl) Metadata(_ context.Context, req ephemeral.MetadataRequest, resp *ephemeral.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_application_v2_client_secret"
}

func (r *resourceImpl) Schema(_ context.Context, _ ephemeral.SchemaRequest, resp *ephemeral.SchemaResponse) {
	resp.Schema = ephschema.Schema{
		MarkdownDescription: "Regenerates and returns the client secret of a `zitadel_application_v2` OIDC or API application " +
			"without persisting it to Terraform state. Evaluating this ephemeral resource rotates the secret, so gate it " +
			"behind `count`/`for_each` to rotate only on demand.",
		Attributes: map[string]ephschema.Attribute{
			"application_id": ephschema.StringAttribute{
				Required:    true,
				Description: "ID of the application whose client secret should be regenerated.",
			},
			"project_id": ephschema.StringAttribute{
				Required:    true,
				Description: "ID of the project the application belongs to.",
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

	client, err := helper.GetAppV2Client(ctx, r.ClientInfo)
	if err != nil {
		resp.Diagnostics.AddError("failed to get client", err.Error())
		return
	}

	// Only override the org from config when supplied; otherwise the call is
	// scoped to the authenticated user's organization.
	if orgID := data.OrgID.ValueString(); orgID != "" {
		ctx = helper.CtxSetOrgID(ctx, orgID)
	}

	zResp, err := client.GenerateClientSecret(ctx, &appv2pb.GenerateClientSecretRequest{
		ApplicationId: data.ApplicationID.ValueString(),
		ProjectId:     data.ProjectID.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("failed to regenerate client secret", err.Error())
		return
	}

	data.ClientSecret = types.StringValue(zResp.GetClientSecret())
	resp.Diagnostics.Append(resp.Result.Set(ctx, &data)...)
}
