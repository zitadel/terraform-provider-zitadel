package zitadel

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	v2 "github.com/zitadel/terraform-provider-zitadel/zitadel/v2"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		DataSourcesMap: map[string]*schema.Resource{
			//		"zitadelV1Org": v1.GetOrgDatasource(),
		},
		Schema: map[string]*schema.Schema{
			v2.DomainVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Domain used to connect to the ZITADEL instance",
			},
			v2.InsecureVar: {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Use insecure connection",
			},
			v2.ProjectVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the ZITADEL project on your ZITADEL instance, for the audience-scope",
			},
			v2.TokenVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Path to the file containing credentials to connect to ZITADEL",
			},
			v2.PortVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used port if not the default ports 80 or 443 are configured",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"zitadel_org":                        v2.OrgResource(),
			"zitadel_human_user":                 v2.GetHumanUser(),
			"zitadel_machine_user":               v2.GetMachineUser(),
			"zitadel_project":                    v2.GetProject(),
			"zitadel_project_role":               v2.GetProjectRole(),
			"zitadel_domain":                     v2.GetDomain(),
			"zitadel_action":                     v2.GetAction(),
			"zitadel_application_oidc":           v2.GetApplicationOIDC(),
			"zitadel_application_api":            v2.GetApplicationAPI(),
			"zitadel_project_grant":              v2.GetProjectGrant(),
			"zitadel_user_grant":                 v2.GetUserGrant(),
			"zitadel_org_member":                 v2.GetOrgMember(),
			"zitadel_project_member":             v2.GetProjectMember(),
			"zitadel_project_grant_member":       v2.GetProjectGrantMember(),
			"zitadel_domain_policy":              v2.GetDomainPolicy(),
			"zitadel_label_policy":               v2.GetLabelPolicy(),
			"zitadel_lockout_policy":             v2.GetLockoutPolicy(),
			"zitadel_login_policy":               v2.GetLoginPolicy(),
			"zitadel_password_complexity_policy": v2.GetPasswordComplexityPolicy(),
			"zitadel_privacy_policy":             v2.GetPrivacyPolicy(),
			"zitadel_trigger_actions":            v2.GetTriggerActions(),
			"zitadel_personal_access_token":      v2.GetPAT(),
			"zitadel_machine_key":                v2.GetMachineKey(),
			"zitadel_org_jwt_idp":                v2.GetOrgJWTIDP(),
			"zitadel_org_oidc_idp":               v2.GetOrgOIDCIDP(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	clientinfo, err := v2.GetClientInfo(d)
	if err != nil {
		return nil, diag.FromErr(err)
	}

	return clientinfo, nil
}
