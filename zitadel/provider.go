package zitadel

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	fdiag "github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	zitadelgo "github.com/zitadel/zitadel-go/v3/pkg/client/zitadel"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/action"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/action_target"
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
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/email_provider_http"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/email_provider_smtp"
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

var _ provider.Provider = (*providerPV6)(nil)

type providerPV6 struct {
	customOptions []zitadelgo.Option
}

func NewProviderPV6(option ...zitadelgo.Option) provider.Provider {
	return &providerPV6{customOptions: option}
}

type providerModel struct {
	Insecure       types.Bool   `tfsdk:"insecure"`
	Domain         types.String `tfsdk:"domain"`
	Port           types.String `tfsdk:"port"`
	AccessToken    types.String `tfsdk:"access_token"`
	Token          types.String `tfsdk:"token"`
	JWTFile        types.String `tfsdk:"jwt_file"`
	JWTProfileFile types.String `tfsdk:"jwt_profile_file"`
	JWTProfileJSON types.String `tfsdk:"jwt_profile_json"`
}

func (p *providerPV6) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "zitadel"
}

func (p *providerPV6) GetSchema(_ context.Context) (tfsdk.Schema, fdiag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			helper.DomainVar: {
				Type:        types.StringType,
				Required:    true,
				Description: helper.DomainDescription,
			},
			helper.InsecureVar: {
				Type:        types.BoolType,
				Optional:    true,
				Description: helper.InsecureDescription,
			},
			helper.AccessTokenVar: {
				Type:        types.StringType,
				Optional:    true,
				Sensitive:   true,
				Description: helper.AccessTokenDescription,
			},
			helper.TokenVar: {
				Type:        types.StringType,
				Optional:    true,
				Description: helper.TokenDescription,
			},
			helper.JWTFileVar: {
				Type:        types.StringType,
				Optional:    true,
				Description: helper.JWTFileDescription,
			},
			helper.JWTProfileFileVar: {
				Type:        types.StringType,
				Optional:    true,
				Description: helper.JWTProfileFileDescription,
			},
			helper.JWTProfileJSONVar: {
				Type:        types.StringType,
				Optional:    true,
				Description: helper.JWTProfileJSONDescription,
			},
			helper.PortVar: {
				Type:        types.StringType,
				Optional:    true,
				Description: helper.PortDescription,
			},
		},
	}, nil
}

func (p *providerPV6) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config providerModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	info, err := helper.GetClientInfo(ctx,
		config.Insecure.ValueBool(),
		config.Domain.ValueString(),
		config.AccessToken.ValueString(),
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

	resp.DataSourceData = info
	resp.ResourceData = info
}

func (p *providerPV6) DataSources(_ context.Context) []func() datasource.DataSource {
	return nil
}

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

func Provider() *schema.Provider {
	return &schema.Provider{
		DataSourcesMap: map[string]*schema.Resource{
			"zitadel_zitadel":                    GetZitadelDatasource(),
			"zitadel_org":                        org.GetDatasource(),
			"zitadel_orgs":                       org.ListDatasources(),
			"zitadel_human_user":                 human_user.GetDatasource(),
			"zitadel_machine_user":               machine_user.GetDatasource(),
			"zitadel_machine_users":              machine_user.ListDatasources(),
			"zitadel_project":                    project.GetDatasource(),
			"zitadel_projects":                   project.ListDatasources(),
			"zitadel_project_role":               project_role.GetDatasource(),
			"zitadel_project_roles":              project_role.ListDatasources(),
			"zitadel_action":                     action.GetDatasource(),
			"zitadel_action_target":              action_target.GetDatasource(),
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
		Schema: map[string]*schema.Schema{
			helper.DomainVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: helper.DomainDescription,
			},
			helper.InsecureVar: {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: helper.InsecureDescription,
			},
			helper.AccessTokenVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: helper.AccessTokenDescription,
				ConflictsWith: []string{
					helper.TokenVar,
					helper.JWTFileVar,
					helper.JWTProfileFileVar,
					helper.JWTProfileJSONVar,
				},
			},
			helper.TokenVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: helper.TokenDescription,
				ConflictsWith: []string{
					helper.AccessTokenVar,
					helper.JWTFileVar,
					helper.JWTProfileFileVar,
					helper.JWTProfileJSONVar,
				},
			},
			helper.JWTFileVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: helper.JWTFileDescription,
				ConflictsWith: []string{
					helper.AccessTokenVar,
					helper.TokenVar,
					helper.JWTProfileFileVar,
					helper.JWTProfileJSONVar,
				},
			},
			helper.JWTProfileFileVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: helper.JWTProfileFileDescription,
				ConflictsWith: []string{
					helper.AccessTokenVar,
					helper.TokenVar,
					helper.JWTFileVar,
					helper.JWTProfileJSONVar,
				},
			},
			helper.JWTProfileJSONVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: helper.JWTProfileJSONDescription,
				ConflictsWith: []string{
					helper.AccessTokenVar,
					helper.TokenVar,
					helper.JWTFileVar,
					helper.JWTProfileFileVar,
				},
			},
			helper.PortVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: helper.PortDescription,
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"zitadel_org":                                org.GetResource(),
			"zitadel_human_user":                         human_user.GetResource(),
			"zitadel_machine_user":                       machine_user.GetResource(),
			"zitadel_project":                            project.GetResource(),
			"zitadel_project_role":                       project_role.GetResource(),
			"zitadel_domain":                             domain.GetResource(),
			"zitadel_action":                             action.GetResource(),
			"zitadel_action_target":                      action_target.GetResource(),
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
			"zitadel_email_provider_smtp":                email_provider_smtp.GetResource(),
			"zitadel_email_provider_http":                email_provider_http.GetResource(),
		},
		ConfigureContextFunc: ProviderConfigure,
	}
}

func ProviderConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	credentials := 0
	for _, k := range []string{
		helper.AccessTokenVar,
		helper.TokenVar,
		helper.JWTFileVar,
		helper.JWTProfileFileVar,
		helper.JWTProfileJSONVar,
	} {
		if v, ok := d.GetOk(k); ok {
			if s, ok := v.(string); ok && s != "" {
				credentials++
			}
		}
	}
	if credentials == 0 {
		return nil, diag.Errorf("one authentication method must be configured")
	}
	if credentials > 1 {
		return nil, diag.Errorf("only one authentication method may be configured")
	}

	clientinfo, err := helper.GetClientInfo(ctx,
		d.Get(helper.InsecureVar).(bool),
		d.Get(helper.DomainVar).(string),
		d.Get(helper.AccessTokenVar).(string),
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
