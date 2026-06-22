// Package organization_domain_validation implements the
// zitadel_organization_domain_validation ephemeral resource.
//
// It (re)generates the verification token for an organization domain and
// returns the token/url for the duration of a single apply only, never writing
// the token to Terraform state (issue #413). Each evaluation generates a fresh
// validation challenge, so gate it behind count/for_each to (re)issue only on
// demand.
package organization_domain_validation

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	ephschema "github.com/hashicorp/terraform-plugin-framework/ephemeral/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/management"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/org"

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
	Domain         types.String `tfsdk:"domain"`
	OrgID          types.String `tfsdk:"org_id"`
	ValidationType types.String `tfsdk:"validation_type"`
	Token          types.String `tfsdk:"token"`
	URL            types.String `tfsdk:"url"`
}

func (r *resourceImpl) Metadata(_ context.Context, req ephemeral.MetadataRequest, resp *ephemeral.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organization_domain_validation"
}

func (r *resourceImpl) Schema(_ context.Context, _ ephemeral.SchemaRequest, resp *ephemeral.SchemaResponse) {
	resp.Schema = ephschema.Schema{
		MarkdownDescription: "Generates and returns the verification token for an organization domain without persisting " +
			"it to Terraform state. Each evaluation generates a fresh validation challenge.",
		Attributes: map[string]ephschema.Attribute{
			"domain": ephschema.StringAttribute{
				Required:    true,
				Description: "The organization domain to generate a validation token for.",
			},
			"org_id": ephschema.StringAttribute{
				Optional:    true,
				Description: "ID of the organization that owns the domain. Defaults to the organization of the authenticated user.",
			},
			"validation_type": ephschema.StringAttribute{
				Required:    true,
				Description: "Type of domain validation" + helper.DescriptionEnumValuesList(org.DomainValidationType_name),
			},
			"token": ephschema.StringAttribute{
				Computed:    true,
				Sensitive:   true,
				Description: "The generated validation token. Only available during apply; never stored in state.",
			},
			"url": ephschema.StringAttribute{
				Computed:    true,
				Description: "The URL at which the token must be served (for HTTP validation).",
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

	// Map the enum name (e.g. DOMAIN_VALIDATION_TYPE_HTTP) to its numeric
	// value, matching the managed zitadel_organization_domain resource. Reject
	// an unknown name loudly rather than silently sending UNSPECIFIED.
	validationType := data.ValidationType.ValueString()
	typeVal, ok := org.DomainValidationType_value[validationType]
	if !ok {
		resp.Diagnostics.AddError("invalid validation_type", fmt.Sprintf("unknown validation type %q", validationType))
		return
	}

	zResp, err := client.GenerateOrgDomainValidation(ctx, &management.GenerateOrgDomainValidationRequest{
		Domain: data.Domain.ValueString(),
		Type:   org.DomainValidationType(typeVal),
	})
	if err != nil {
		resp.Diagnostics.AddError("failed to generate org domain validation", err.Error())
		return
	}

	data.Token = types.StringValue(zResp.GetToken())
	data.URL = types.StringValue(zResp.GetUrl())
	resp.Diagnostics.Append(resp.Result.Set(ctx, &data)...)
}
