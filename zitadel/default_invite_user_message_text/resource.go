package default_invite_user_message_text

import (
	"context"
	"fmt"

	github_com_hashicorp_terraform_plugin_framework_diag "github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	github_com_hashicorp_terraform_plugin_framework_tfsdk "github.com/hashicorp/terraform-plugin-framework/tfsdk"
	github_com_hashicorp_terraform_plugin_framework_types "github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/admin"
)

const (
	LanguageVar = "language"
)

var (
	_ resource.Resource = &defaultInviteUserMessageTextResource{}
)

type defaultInviteUserMessageTextModel struct {
	ID         github_com_hashicorp_terraform_plugin_framework_types.String `tfsdk:"id"`
	Language   github_com_hashicorp_terraform_plugin_framework_types.String `tfsdk:"language"`
	Title      github_com_hashicorp_terraform_plugin_framework_types.String `tfsdk:"title"`
	PreHeader  github_com_hashicorp_terraform_plugin_framework_types.String `tfsdk:"pre_header"`
	Subject    github_com_hashicorp_terraform_plugin_framework_types.String `tfsdk:"subject"`
	Greeting   github_com_hashicorp_terraform_plugin_framework_types.String `tfsdk:"greeting"`
	Text       github_com_hashicorp_terraform_plugin_framework_types.String `tfsdk:"text"`
	ButtonText github_com_hashicorp_terraform_plugin_framework_types.String `tfsdk:"button_text"`
	FooterText github_com_hashicorp_terraform_plugin_framework_types.String `tfsdk:"footer_text"`
}

func New() resource.Resource {
	return &defaultInviteUserMessageTextResource{}
}

type defaultInviteUserMessageTextResource struct {
	clientInfo *helper.ClientInfo
}

func (r *defaultInviteUserMessageTextResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_invite_user_message_text"
}

func (r *defaultInviteUserMessageTextResource) GetSchema(_ context.Context) (github_com_hashicorp_terraform_plugin_framework_tfsdk.Schema, github_com_hashicorp_terraform_plugin_framework_diag.Diagnostics) {
	return github_com_hashicorp_terraform_plugin_framework_tfsdk.Schema{
		Description: "Resource for managing the default invite user message texts within ZITADEL.",
		Attributes: map[string]github_com_hashicorp_terraform_plugin_framework_tfsdk.Attribute{
			"id": {
				Computed:    true,
				Optional:    false,
				Required:    false,
				Type:        github_com_hashicorp_terraform_plugin_framework_types.StringType,
				Description: "The ID of this resource. Equal to `language`.",
				PlanModifiers: []github_com_hashicorp_terraform_plugin_framework_tfsdk.AttributePlanModifier{
					resource.UseStateForUnknown(),
				},
			},
			"language": {
				Required:    true, // language remains required, as it's the ID
				Optional:    false,
				Computed:    false,
				Type:        github_com_hashicorp_terraform_plugin_framework_types.StringType,
				Description: "The language of the invite user message text (e.g., `en`, `de`). This also serves as the resource's ID.",
			},
			"title": {
				Optional:    true, // Changed to Optional
				Required:    false,
				Computed:    false,
				Type:        github_com_hashicorp_terraform_plugin_framework_types.StringType,
				Description: "The title of the invite user message.",
			},
			"pre_header": {
				Optional:    true, // Changed to Optional
				Required:    false,
				Computed:    false,
				Type:        github_com_hashicorp_terraform_plugin_framework_types.StringType,
				Description: "The pre-header text of the invite user message.",
			},
			"subject": {
				Optional:    true, // Changed to Optional
				Required:    false,
				Computed:    false,
				Type:        github_com_hashicorp_terraform_plugin_framework_types.StringType,
				Description: "The subject line of the invite user message.",
			},
			"greeting": {
				Optional:    true, // Changed to Optional
				Required:    false,
				Computed:    false,
				Type:        github_com_hashicorp_terraform_plugin_framework_types.StringType,
				Description: "The greeting text of the invite user message.",
			},
			"text": {
				Optional:    true, // Changed to Optional
				Required:    false,
				Computed:    false,
				Type:        github_com_hashicorp_terraform_plugin_framework_types.StringType,
				Description: "The main body text of the invite user message.",
			},
			"button_text": {
				Optional:    true, // Changed to Optional
				Required:    false,
				Computed:    false,
				Type:        github_com_hashicorp_terraform_plugin_framework_types.StringType,
				Description: "The text displayed on the call-to-action button in the message.",
			},
			"footer_text": {
				Optional:    true, // Changed to Optional
				Required:    false,
				Computed:    false,
				Type:        github_com_hashicorp_terraform_plugin_framework_types.StringType,
				Description: "The footer text of the invite user message.",
			},
		},
	}, nil
}

func (r *defaultInviteUserMessageTextResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.clientInfo = req.ProviderData.(*helper.ClientInfo)
}

func (r *defaultInviteUserMessageTextResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan defaultInviteUserMessageTextModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := helper.GetAdminClient(ctx, r.clientInfo)
	if err != nil {
		resp.Diagnostics.AddError("failed to get client", err.Error())
		return
	}

	zReq := &admin.SetDefaultInviteUserMessageTextRequest{
		Language:   plan.Language.ValueString(),
		Title:      plan.Title.ValueString(),
		PreHeader:  plan.PreHeader.ValueString(),
		Subject:    plan.Subject.ValueString(),
		Greeting:   plan.Greeting.ValueString(),
		Text:       plan.Text.ValueString(),
		ButtonText: plan.ButtonText.ValueString(),
		FooterText: plan.FooterText.ValueString(),
	}

	_, err = client.SetDefaultInviteUserMessageText(ctx, zReq)
	if err != nil {
		resp.Diagnostics.AddError("failed to create default invite user message text", fmt.Sprintf("Error: %s", err.Error()))
		return
	}

	plan.ID = plan.Language

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *defaultInviteUserMessageTextResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state defaultInviteUserMessageTextModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	language := state.ID.ValueString()

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

	state.Title = github_com_hashicorp_terraform_plugin_framework_types.StringValue(zResp.CustomText.Title)
	state.PreHeader = github_com_hashicorp_terraform_plugin_framework_types.StringValue(zResp.CustomText.PreHeader)
	state.Subject = github_com_hashicorp_terraform_plugin_framework_types.StringValue(zResp.CustomText.Subject)
	state.Greeting = github_com_hashicorp_terraform_plugin_framework_types.StringValue(zResp.CustomText.Greeting)
	state.Text = github_com_hashicorp_terraform_plugin_framework_types.StringValue(zResp.CustomText.Text)
	state.ButtonText = github_com_hashicorp_terraform_plugin_framework_types.StringValue(zResp.CustomText.ButtonText)
	state.FooterText = github_com_hashicorp_terraform_plugin_framework_types.StringValue(zResp.CustomText.FooterText)
	state.ID = github_com_hashicorp_terraform_plugin_framework_types.StringValue(language)
	state.Language = github_com_hashicorp_terraform_plugin_framework_types.StringValue(language)

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
	if resp.Diagnostics.HasError() { // Re-added this check for robustness
		return
	}
}

func (r *defaultInviteUserMessageTextResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan defaultInviteUserMessageTextModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := helper.GetAdminClient(ctx, r.clientInfo)
	if err != nil {
		resp.Diagnostics.AddError("Failed to get client", err.Error())
		return
	}

	zReq := &admin.SetDefaultInviteUserMessageTextRequest{
		Language:   plan.Language.ValueString(),
		Title:      plan.Title.ValueString(),
		PreHeader:  plan.PreHeader.ValueString(),
		Subject:    plan.Subject.ValueString(),
		Greeting:   plan.Greeting.ValueString(),
		Text:       plan.Text.ValueString(),
		ButtonText: plan.ButtonText.ValueString(),
		FooterText: plan.FooterText.ValueString(),
	}

	_, err = client.SetDefaultInviteUserMessageText(ctx, zReq)
	if err != nil {
		resp.Diagnostics.AddError("Failed to update default invite user message text", err.Error())
		return
	}

	// No explicit setID(&plan, ...) call needed here.
	// The `plan` struct already contains the desired ID and Language from the input.

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *defaultInviteUserMessageTextResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var model defaultInviteUserMessageTextModel
	resp.Diagnostics.Append(req.State.Get(ctx, &model)...)
	if resp.Diagnostics.HasError() {
		return
	}

	language := model.ID.ValueString()
	if language == "" {
		// If language/ID is empty, it means the resource wasn't properly in state.
		// The framework will remove it if this method completes without error.
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

	// In Terraform Plugin Framework v0.15.0, if the Delete method completes
	// without returning an error or adding errors to resp.Diagnostics,
	// the framework automatically removes the resource from the Terraform state.
}

func setID(model *defaultInviteUserMessageTextModel, language string) {
	model.ID = github_com_hashicorp_terraform_plugin_framework_types.StringValue(language)
	model.Language = github_com_hashicorp_terraform_plugin_framework_types.StringValue(language)
}
