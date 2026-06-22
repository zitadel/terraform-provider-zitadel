// Package personal_access_token implements the zitadel_personal_access_token
// ephemeral resource.
//
// It creates a new personal access token (PAT) for a machine (service) user and
// returns the token for the duration of a single apply only, never writing it
// to Terraform state (issue #413). Each evaluation mints a new token, so gate
// it behind count/for_each to issue a token only on demand. Existing tokens
// remain valid until removed via the managed zitadel_personal_access_token
// resource or the API.
package personal_access_token

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	ephschema "github.com/hashicorp/terraform-plugin-framework/ephemeral/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/management"
	"google.golang.org/protobuf/types/known/timestamppb"

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
	UserID         types.String `tfsdk:"user_id"`
	OrgID          types.String `tfsdk:"org_id"`
	ExpirationDate types.String `tfsdk:"expiration_date"`
	TokenID        types.String `tfsdk:"token_id"`
	Token          types.String `tfsdk:"token"`
}

func (r *resourceImpl) Metadata(_ context.Context, req ephemeral.MetadataRequest, resp *ephemeral.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_personal_access_token"
}

func (r *resourceImpl) Schema(_ context.Context, _ ephemeral.SchemaRequest, resp *ephemeral.SchemaResponse) {
	resp.Schema = ephschema.Schema{
		MarkdownDescription: "Creates a new personal access token for a machine (service) user and returns it without " +
			"persisting it to Terraform state. Each evaluation mints a new token, so gate it behind `count`/`for_each`.",
		Attributes: map[string]ephschema.Attribute{
			"user_id": ephschema.StringAttribute{
				Required:    true,
				Description: "ID of the machine user the token is created for.",
			},
			"org_id": ephschema.StringAttribute{
				Optional:    true,
				Description: "ID of the organization that owns the machine user. Defaults to the organization of the authenticated user.",
			},
			"expiration_date": ephschema.StringAttribute{
				Optional:    true,
				Description: "Optional expiration date for the token, in RFC3339 format (e.g. 2519-04-01T08:45:00Z).",
			},
			"token_id": ephschema.StringAttribute{
				Computed:    true,
				Description: "ID of the generated token.",
			},
			"token": ephschema.StringAttribute{
				Computed:    true,
				Sensitive:   true,
				Description: "The generated personal access token. Only available during apply; never stored in state.",
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

	zReq := &management.AddPersonalAccessTokenRequest{
		UserId: data.UserID.ValueString(),
	}
	if exp := data.ExpirationDate.ValueString(); exp != "" {
		t, perr := time.Parse(time.RFC3339, exp)
		if perr != nil {
			resp.Diagnostics.AddError("failed to parse expiration_date", perr.Error())
			return
		}
		zReq.ExpirationDate = timestamppb.New(t)
	}

	zResp, err := client.AddPersonalAccessToken(ctx, zReq)
	if err != nil {
		resp.Diagnostics.AddError("failed to add personal access token", err.Error())
		return
	}

	data.TokenID = types.StringValue(zResp.GetTokenId())
	data.Token = types.StringValue(zResp.GetToken())
	resp.Diagnostics.Append(resp.Result.Set(ctx, &data)...)
}
