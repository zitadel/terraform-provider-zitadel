// Package application_key implements the zitadel_application_key ephemeral
// resource.
//
// It creates a new (JSON) key for an API application and returns the
// key_details for the duration of a single apply only, never writing them to
// Terraform state (issue #413). Each evaluation adds a new key, so gate it
// behind count/for_each to add a key only on demand. Old keys remain valid
// until removed via the managed zitadel_application_key resource or the API.
package application_key

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	ephschema "github.com/hashicorp/terraform-plugin-framework/ephemeral/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/authn"
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
	ProjectID      types.String `tfsdk:"project_id"`
	AppID          types.String `tfsdk:"app_id"`
	OrgID          types.String `tfsdk:"org_id"`
	ExpirationDate types.String `tfsdk:"expiration_date"`
	KeyID          types.String `tfsdk:"key_id"`
	KeyDetails     types.String `tfsdk:"key_details"`
}

func (r *resourceImpl) Metadata(_ context.Context, req ephemeral.MetadataRequest, resp *ephemeral.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_application_key"
}

func (r *resourceImpl) Schema(_ context.Context, _ ephemeral.SchemaRequest, resp *ephemeral.SchemaResponse) {
	resp.Schema = ephschema.Schema{
		MarkdownDescription: "Creates a new JSON key for an API application and returns the key material without persisting " +
			"it to Terraform state. Each evaluation adds a new key, so gate it behind `count`/`for_each`.",
		Attributes: map[string]ephschema.Attribute{
			"project_id": ephschema.StringAttribute{
				Required:    true,
				Description: "ID of the project the application belongs to.",
			},
			"app_id": ephschema.StringAttribute{
				Required:    true,
				Description: "ID of the application the key is created for.",
			},
			"org_id": ephschema.StringAttribute{
				Optional:    true,
				Description: "ID of the organization that owns the application. Defaults to the organization of the authenticated user.",
			},
			"expiration_date": ephschema.StringAttribute{
				Optional:    true,
				Description: "Optional expiration date for the key, in RFC3339 format (e.g. 2519-04-01T08:45:00Z).",
			},
			"key_id": ephschema.StringAttribute{
				Computed:    true,
				Description: "ID of the generated key.",
			},
			"key_details": ephschema.StringAttribute{
				Computed:    true,
				Sensitive:   true,
				Description: "The generated key material (JSON). Only available during apply; never stored in state.",
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

	zReq := &management.AddAppKeyRequest{
		ProjectId: data.ProjectID.ValueString(),
		AppId:     data.AppID.ValueString(),
		// Only JSON keys are supported, matching the managed resource.
		Type: authn.KeyType_KEY_TYPE_JSON,
	}
	if exp := data.ExpirationDate.ValueString(); exp != "" {
		t, perr := time.Parse(time.RFC3339, exp)
		if perr != nil {
			resp.Diagnostics.AddError("failed to parse expiration_date", perr.Error())
			return
		}
		zReq.ExpirationDate = timestamppb.New(t)
	}

	zResp, err := client.AddAppKey(ctx, zReq)
	if err != nil {
		resp.Diagnostics.AddError("failed to add application key", err.Error())
		return
	}

	data.KeyID = types.StringValue(zResp.GetId())
	data.KeyDetails = types.StringValue(string(zResp.GetKeyDetails()))
	resp.Diagnostics.Append(resp.Result.Set(ctx, &data)...)
}
