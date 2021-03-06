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
			v2.IssuerVar: {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("ISSUER", ""),
			},
			v2.AddressVar: {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("ADDRESS", ""),
			},
			v2.ProjectVar: {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("PROJECT", ""),
			},
			v2.TokenVar: {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("TOKEN", ""),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"org":                  v2.OrgResource(),
			"human_user":           v2.GetHumanUser(),
			"machine_user":         v2.GetMachineUser(),
			"project":              v2.GetProject(),
			"project_role":         v2.GetProjectRole(),
			"domain":               v2.GetDomain(),
			"action":               v2.GetAction(),
			"application_oidc":     v2.GetApplicationOIDC(),
			"application_api":      v2.GetApplicationAPI(),
			"project_grant":        v2.GetProjectGrant(),
			"user_grant":           v2.GetUserGrant(),
			"org_member":           v2.GetOrgMember(),
			"project_member":       v2.GetProjectMember(),
			"project_grant_member": v2.GetProjectGrantMember(),
			/*
				"domain_policy":              v2.GetDomainPolicy(),
				"label_policy":               v2.GetLabelPolicy(),
				"lockout_policy":             v2.GetLockoutPolicy(),
				"login_policy":               v2.GetLoginPolicy(),
				"password_complexity_policy": v2.GetPasswordComplexityPolicy(),
				"privacy_policy":             v2.GetPrivacyPolicy(),
			*/
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
