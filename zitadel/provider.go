package zitadel

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	sdkschema "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	zitadel_go "github.com/zitadel/zitadel-go/v3/pkg/client/zitadel"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/action"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/application_api"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/application_key"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/application_oidc"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/application_saml"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/default_domain_claimed_message_text"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/default_domain_policy"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/default_init_message_text"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/default_label_policy"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/default_lockout_policy"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/default_login_policy"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/default_login_texts"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/default_notification_policy"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/default_oidc_settings"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/default_password_age_policy"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/default_password_change_message_text"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/default_password_complexity_policy"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/default_password_reset_message_text"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/default_passwordless_registration_message_text"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/default_privacy_policy"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/default_verify_email_message_text"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/default_verify_email_otp_message_text"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/default_verify_phone_message_text"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/default_verify_sms_otp_message_text"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/domain"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/domain_claimed_message_text"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/domain_policy"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/human_user"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/idp_azure_ad"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/idp_github"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/idp_github_es"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/idp_gitlab"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/idp_gitlab_self_hosted"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/idp_google"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/idp_ldap"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/idp_oauth"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/idp_oidc"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/idp_saml"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/init_message_text"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/instance_member"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/label_policy"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/lockout_policy"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/login_policy"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/login_texts"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/machine_key"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/machine_user"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/notification_policy"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/org"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/org_idp_azure_ad"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/org_idp_github"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/org_idp_github_es"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/org_idp_gitlab"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/org_idp_gitlab_self_hosted"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/org_idp_google"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/org_idp_jwt"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/org_idp_ldap"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/org_idp_oauth"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/org_idp_oidc"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/org_idp_saml"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/org_member"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/org_metadata"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/password_age_policy"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/password_change_message_text"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/password_complexity_policy"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/password_reset_message_text"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/passwordless_registration_message_text"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/pat"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/privacy_policy"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/project"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/project_grant"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/project_grant_member"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/project_member"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/project_role"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/sms_provider_http"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/sms_provider_twilio"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/smtp_config"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/trigger_actions"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/user_grant"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/user_metadata"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/verify_email_message_text"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/verify_email_otp_message_text"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/verify_phone_message_text"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/verify_sms_otp_message_text"
)

// Ensure provider satisfies various provider interfaces
var _ provider.Provider = (*providerPV6)(nil)

// providerPV6 is the provider implementation for the Terraform Plugin Framework v6
type providerPV6 struct {
	customOptions []zitadel_go.Option
}

// NewProviderPV6 creates a new provider instance with optional configuration
func NewProviderPV6(option ...zitadel_go.Option) provider.Provider {
	return &providerPV6{customOptions: option}
}

// providerModel represents the provider's configuration schema
type providerModel struct {
	Insecure       types.Bool   `tfsdk:"insecure"`
	Domain         types.String `tfsdk:"domain"`
	Port           types.String `tfsdk:"port"`
	Token          types.String `tfsdk:"token"`
	JWTFile        types.String `tfsdk:"jwt_file"`
	JWTProfileFile types.String `tfsdk:"jwt_profile_file"`
	JWTProfileJSON types.String `tfsdk:"jwt_profile_json"`
}

// Metadata returns the provider type name
func (p *providerPV6) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "zitadel"
}

// Schema defines the provider-level schema for configuration data
func (p *providerPV6) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			helper.DomainVar: schema.StringAttribute{
				Required:    true,
				Description: helper.DomainDescription,
			},
			helper.InsecureVar: schema.BoolAttribute{
				Optional:    true,
				Description: helper.InsecureDescription,
			},
			helper.TokenVar: schema.StringAttribute{
				Optional:    true,
				Description: helper.TokenDescription,
			},
			helper.JWTFileVar: schema.StringAttribute{
				Optional:    true,
				Description: helper.JWTFileDescription,
			},
			helper.JWTProfileFileVar: schema.StringAttribute{
				Optional:    true,
				Description: helper.JWTProfileFileDescription,
			},
			helper.JWTProfileJSONVar: schema.StringAttribute{
				Optional:    true,
				Description: helper.JWTProfileJSONDescription,
			},
			helper.PortVar: schema.StringAttribute{
				Optional:    true,
				Description: helper.PortDescription,
			},
		},
	}
}

// Configure prepares the provider for data source and resource CRUD operations
func (p *providerPV6) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config providerModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Initialize the client with configuration values
	info, err := helper.GetClientInfo(ctx,
		config.Insecure.ValueBool(),
		config.Domain.ValueString(),
		config.Token.ValueString(),
		config.JWTFile.ValueString(),
		config.JWTProfileFile.ValueString(),
		config.JWTProfileJSON.ValueString(),
		config.Port.ValueString(),
	)
	if err != nil {
		resp.Diagnostics.AddError("failed to handle provider config", err.Error())
		return
	}

	// Make the client configuration available to resources and data sources
	resp.DataSourceData = info
	resp.ResourceData = info
}

// DataSources defines the data sources implemented in the provider
func (p *providerPV6) DataSources(_ context.Context) []func() datasource.DataSource {
	return nil
}

// Resources defines the resources implemented in the provider
func (p *providerPV6) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		init_message_text.New,
		login_texts.New,
		password_reset_message_text.New,
		password_change_message_text.New,
		verify_email_message_text.New,
		verify_phone_message_text.New,
		domain_claimed_message_text.New,
		passwordless_registration_message_text.New,
		default_domain_claimed_message_text.New,
		default_init_message_text.New,
		default_login_texts.New,
		default_password_reset_message_text.New,
		default_password_change_message_text.New,
		default_passwordless_registration_message_text.New,
		default_verify_email_message_text.New,
		default_verify_phone_message_text.New,
		default_verify_email_otp_message_text.New,
		default_verify_sms_otp_message_text.New,
		verify_email_otp_message_text.New,
		verify_sms_otp_message_text.New,
	}
}

// Provider returns the SDK v2 provider for backward compatibility
// This maintains support for existing configurations while transitioning to Framework v6
func Provider() *sdkschema.Provider {
	return &sdkschema.Provider{
		DataSourcesMap: map[string]*sdkschema.Resource{
			"zitadel_org":                        org.GetDatasource(),
			"zitadel_orgs":                       org.ListDatasources(),
			"zitadel_human_user":                 human_user.GetDatasource(),
			"zitadel_machine_user":               machine_user.GetDatasource(),
			"zitadel_machine_users":              machine_user.ListDatasources(),
			"zitadel_project":                    project.GetDatasource(),
			"zitadel_projects":                   project.ListDatasources(),
			"zitadel_project_role":               project_role.GetDatasource(),
			"zitadel_action":                     action.GetDatasource(),
			"zitadel_application_oidc":           application_oidc.GetDatasource(),
			"zitadel_application_oidcs":          application_oidc.ListDatasources(),
			"zitadel_application_api":            application_api.GetDatasource(),
			"zitadel_application_apis":           application_api.ListDatasources(),
			"zitadel_application_saml":           application_saml.GetDatasource(),
			"zitadel_application_samls":          application_saml.ListDatasources(),
			"zitadel_trigger_actions":            trigger_actions.GetDatasource(),
			"zitadel_idp_github":                 idp_github.GetDatasource(),
			"zitadel_idp_github_es":              idp_github_es.GetDatasource(),
			"zitadel_idp_gitlab":                 idp_gitlab.GetDatasource(),
			"zitadel_idp_gitlab_self_hosted":     idp_gitlab_self_hosted.GetDatasource(),
			"zitadel_idp_google":                 idp_google.GetDatasource(),
			"zitadel_idp_azure_ad":               idp_azure_ad.GetDatasource(),
			"zitadel_idp_ldap":                   idp_ldap.GetDatasource(),
			"zitadel_idp_saml":                   idp_saml.GetDatasource(),
			"zitadel_idp_oauth":                  idp_oauth.GetDatasource(),
			"zitadel_idp_oidc":                   idp_oidc.GetDatasource(),
			"zitadel_org_jwt_idp":                org_idp_jwt.GetDatasource(),
			"zitadel_org_oidc_idp":               org_idp_oidc.GetDatasource(),
			"zitadel_org_idp_github":             org_idp_github.GetDatasource(),
			"zitadel_org_idp_github_es":          org_idp_github_es.GetDatasource(),
			"zitadel_org_idp_gitlab":             org_idp_gitlab.GetDatasource(),
			"zitadel_org_idp_gitlab_self_hosted": org_idp_gitlab_self_hosted.GetDatasource(),
			"zitadel_org_idp_google":             org_idp_google.GetDatasource(),
			"zitadel_org_idp_azure_ad":           org_idp_azure_ad.GetDatasource(),
			"zitadel_org_idp_ldap":               org_idp_ldap.GetDatasource(),
			"zitadel_org_idp_saml":               org_idp_saml.GetDatasource(),
			"zitadel_org_idp_oauth":              org_idp_oauth.GetDatasource(),
			"zitadel_default_oidc_settings":      default_oidc_settings.GetDatasource(),
		},
		Schema: map[string]*sdkschema.Schema{
			helper.DomainVar: {
				Type:        sdkschema.TypeString,
				Required:    true,
				Description: helper.DomainDescription,
			},
			helper.InsecureVar: {
				Type:        sdkschema.TypeBool,
				Optional:    true,
				Description: helper.InsecureDescription,
			},
			helper.TokenVar: {
				Type:        sdkschema.TypeString,
				Optional:    true,
				Description: helper.TokenDescription,
			},
			helper.JWTFileVar: {
				Type:        sdkschema.TypeString,
				Optional:    true,
				Description: helper.JWTFileDescription,
			},
			helper.JWTProfileFileVar: {
				Type:        sdkschema.TypeString,
				Optional:    true,
				Description: helper.JWTProfileFileDescription,
			},
			helper.JWTProfileJSONVar: {
				Type:        sdkschema.TypeString,
				Optional:    true,
				Description: helper.JWTProfileJSONDescription,
			},
			helper.PortVar: {
				Type:        sdkschema.TypeString,
				Optional:    true,
				Description: helper.PortDescription,
			},
		},
		ResourcesMap: map[string]*sdkschema.Resource{
			"zitadel_org":                                org.GetResource(),
			"zitadel_human_user":                         human_user.GetResource(),
			"zitadel_machine_user":                       machine_user.GetResource(),
			"zitadel_project":                            project.GetResource(),
			"zitadel_project_role":                       project_role.GetResource(),
			"zitadel_domain":                             domain.GetResource(),
			"zitadel_action":                             action.GetResource(),
			"zitadel_application_oidc":                   application_oidc.GetResource(),
			"zitadel_application_api":                    application_api.GetResource(),
			"zitadel_application_saml":                   application_saml.GetResource(),
			"zitadel_application_key":                    application_key.GetResource(),
			"zitadel_project_grant":                      project_grant.GetResource(),
			"zitadel_user_grant":                         user_grant.GetResource(),
			"zitadel_org_member":                         org_member.GetResource(),
			"zitadel_instance_member":                    instance_member.GetResource(),
			"zitadel_project_member":                     project_member.GetResource(),
			"zitadel_project_grant_member":               project_grant_member.GetResource(),
			"zitadel_domain_policy":                      domain_policy.GetResource(),
			"zitadel_label_policy":                       label_policy.GetResource(),
			"zitadel_lockout_policy":                     lockout_policy.GetResource(),
			"zitadel_login_policy":                       login_policy.GetResource(),
			"zitadel_password_age_policy":                password_age_policy.GetResource(),
			"zitadel_password_complexity_policy":         password_complexity_policy.GetResource(),
			"zitadel_privacy_policy":                     privacy_policy.GetResource(),
			"zitadel_trigger_actions":                    trigger_actions.GetResource(),
			"zitadel_personal_access_token":              pat.GetResource(),
			"zitadel_machine_key":                        machine_key.GetResource(),
			"zitadel_default_label_policy":               default_label_policy.GetResource(),
			"zitadel_default_login_policy":               default_login_policy.GetResource(),
			"zitadel_default_lockout_policy":             default_lockout_policy.GetResource(),
			"zitadel_default_domain_policy":              default_domain_policy.GetResource(),
			"zitadel_default_privacy_policy":             default_privacy_policy.GetResource(),
			"zitadel_default_password_age_policy":        default_password_age_policy.GetResource(),
			"zitadel_default_password_complexity_policy": default_password_complexity_policy.GetResource(),
			"zitadel_sms_provider_twilio":                sms_provider_twilio.GetResource(),
			"zitadel_sms_provider_http":                  sms_provider_http.GetResource(),
			"zitadel_smtp_config":                        smtp_config.GetResource(),
			"zitadel_default_notification_policy":        default_notification_policy.GetResource(),
			"zitadel_notification_policy":                notification_policy.GetResource(),
			"zitadel_idp_github":                         idp_github.GetResource(),
			"zitadel_idp_github_es":                      idp_github_es.GetResource(),
			"zitadel_idp_gitlab":                         idp_gitlab.GetResource(),
			"zitadel_idp_gitlab_self_hosted":             idp_gitlab_self_hosted.GetResource(),
			"zitadel_idp_google":                         idp_google.GetResource(),
			"zitadel_idp_azure_ad":                       idp_azure_ad.GetResource(),
			"zitadel_idp_ldap":                           idp_ldap.GetResource(),
			"zitadel_idp_saml":                           idp_saml.GetResource(),
			"zitadel_idp_oauth":                          idp_oauth.GetResource(),
			"zitadel_idp_oidc":                           idp_oidc.GetResource(),
			"zitadel_org_idp_jwt":                        org_idp_jwt.GetResource(),
			"zitadel_org_idp_oidc":                       org_idp_oidc.GetResource(),
			"zitadel_org_idp_github":                     org_idp_github.GetResource(),
			"zitadel_org_idp_github_es":                  org_idp_github_es.GetResource(),
			"zitadel_org_idp_gitlab":                     org_idp_gitlab.GetResource(),
			"zitadel_org_idp_gitlab_self_hosted":         org_idp_gitlab_self_hosted.GetResource(),
			"zitadel_org_idp_google":                     org_idp_google.GetResource(),
			"zitadel_org_idp_azure_ad":                   org_idp_azure_ad.GetResource(),
			"zitadel_org_idp_ldap":                       org_idp_ldap.GetResource(),
			"zitadel_org_idp_saml":                       org_idp_saml.GetResource(),
			"zitadel_org_idp_oauth":                      org_idp_oauth.GetResource(),
			"zitadel_default_oidc_settings":              default_oidc_settings.GetResource(),
			"zitadel_org_metadata":                       org_metadata.GetResource(),
			"zitadel_user_metadata":                      user_metadata.GetResource(),
		},
		ConfigureContextFunc: ProviderConfigure,
	}
}

// ProviderConfigure configures the SDK v2 provider for backward compatibility
func ProviderConfigure(ctx context.Context, d *sdkschema.ResourceData) (interface{}, diag.Diagnostics) {
	clientinfo, err := helper.GetClientInfo(ctx,
		d.Get(helper.InsecureVar).(bool),
		d.Get(helper.DomainVar).(string),
		d.Get(helper.TokenVar).(string),
		d.Get(helper.JWTFileVar).(string),
		d.Get(helper.JWTProfileFileVar).(string),
		d.Get(helper.JWTProfileJSONVar).(string),
		d.Get(helper.PortVar).(string),
	)
	if err != nil {
		return nil, diag.FromErr(err)
	}

	return clientinfo, nil
}