package zitadel

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/action"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/app_key"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/application_api"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/application_oidc"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/domain"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/domain_policy"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/helper"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/human_user"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/idp_jwt"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/idp_oidc"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/label_policy"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/lockout_policy"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/login_policy"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/machine_key"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/machine_user"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/org"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/org_member"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/password_complexity_policy"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/pat"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/privacy_policy"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/project"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/project_grant"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/project_grant_member"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/project_member"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/project_role"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/smtp_config"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/trigger_actions"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/user_grant"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		DataSourcesMap: map[string]*schema.Resource{
			"zitadel_org":              org.GetDatasource(),
			"zitadel_human_user":       human_user.GetDatasource(),
			"zitadel_machine_user":     machine_user.GetDatasource(),
			"zitadel_project":          project.GetDatasource(),
			"zitadel_project_role":     project_role.GetDatasource(),
			"zitadel_action":           action.GetDatasource(),
			"zitadel_application_oidc": application_oidc.GetDatasource(),
			"zitadel_application_api":  application_api.GetDatasource(),
			"zitadel_trigger_actions":  trigger_actions.GetDatasource(),
			"zitadel_org_jwt_idp":      idp_jwt.GetDatasource(),
			"zitadel_org_oidc_idp":     idp_oidc.GetDatasource(),
		},
		Schema: map[string]*schema.Schema{
			helper.DomainVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Domain used to connect to the ZITADEL instance",
			},
			helper.InsecureVar: {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Use insecure connection",
			},
			helper.TokenVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Path to the file containing credentials to connect to ZITADEL",
			},
			helper.PortVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used port if not the default ports 80 or 443 are configured",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"zitadel_org":                        org.GetResource(),
			"zitadel_human_user":                 human_user.GetResource(),
			"zitadel_machine_user":               machine_user.GetResource(),
			"zitadel_project":                    project.GetResource(),
			"zitadel_project_role":               project_role.GetResource(),
			"zitadel_domain":                     domain.GetResource(),
			"zitadel_action":                     action.GetResource(),
			"zitadel_application_oidc":           application_oidc.GetResource(),
			"zitadel_application_api":            application_api.GetResource(),
			"zitadel_application_key":            app_key.GetResource(),
			"zitadel_project_grant":              project_grant.GetResource(),
			"zitadel_user_grant":                 user_grant.GetResource(),
			"zitadel_org_member":                 org_member.GetResource(),
			"zitadel_project_member":             project_member.GetResource(),
			"zitadel_project_grant_member":       project_grant_member.GetResource(),
			"zitadel_domain_policy":              domain_policy.GetResource(),
			"zitadel_label_policy":               label_policy.GetResource(),
			"zitadel_lockout_policy":             lockout_policy.GetResource(),
			"zitadel_login_policy":               login_policy.GetResource(),
			"zitadel_password_complexity_policy": password_complexity_policy.GetResource(),
			"zitadel_privacy_policy":             privacy_policy.GetResource(),
			"zitadel_trigger_actions":            trigger_actions.GetResource(),
			"zitadel_personal_access_token":      pat.GetResource(),
			"zitadel_machine_key":                machine_key.GetResource(),
			"zitadel_org_jwt_idp":                idp_jwt.GetResource(),
			"zitadel_org_oidc_idp":               idp_oidc.GetResource(),
			"zitadel_smtp_config":                smtp_config.GetResource(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	clientinfo, err := helper.GetClientInfo(d)
	if err != nil {
		return nil, diag.FromErr(err)
	}

	return clientinfo, nil
}
