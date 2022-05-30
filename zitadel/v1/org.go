package v1

import (
	"context"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	management2 "github.com/zitadel/zitadel-go/pkg/client/zitadel/management"
)

const (
	orgVar     = "org"
	nameVar    = "name"
	issuerVar  = "issuer"
	addressVar = "address"
	projectVar = "project"
	tokenVar   = "token"

	passwordComplexityPolicyVar = "password_complexity_policy"
	lockoutPolicyVar            = "lockout_policy"
	loginPolicyVar              = "login_policy"
	iamPolicyVar                = "iam_policy"
	labelPolicyVar              = "label_policy"
	privacyPolicyVar            = "privacy_policy"

	usersVar    = "users"
	projectsVar = "projects"
	domainsVar  = "domains"
	actionsVar  = "actions"
)

func GetOrgDatasource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			orgVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID of the organization",
			},
			nameVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the organization",
			},
			issuerVar: {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("ISSUER", ""),
			},
			addressVar: {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("ADDRESS", ""),
			},
			projectVar: {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("PROJECT", ""),
			},
			tokenVar: {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("SERVICE_TOKEN", ""),
			},

			passwordComplexityPolicyVar: {
				Type:        schema.TypeSet,
				Elem:        GetPasswordComplexityPolicyDatasource(),
				Optional:    true,
				Computed:    true,
				Description: "List of password complexity policies in organization",
			},
			lockoutPolicyVar: {
				Type:        schema.TypeSet,
				Elem:        GetLockoutPolicyDatasource(),
				Optional:    true,
				Computed:    true,
				Description: "List of lockout policies in organization",
			},
			loginPolicyVar: {
				Type:        schema.TypeSet,
				Elem:        GetLoginPolicyDatasource(),
				Optional:    true,
				Computed:    true,
				Description: "List of login policies in organization",
			},
			labelPolicyVar: {
				Type:        schema.TypeSet,
				Elem:        GetLabelPolicyDatasource(),
				Optional:    true,
				Computed:    true,
				Description: "List of label policies in organization",
			},
			iamPolicyVar: {
				Type:        schema.TypeSet,
				Elem:        GetIAMPolicyDatasource(),
				Optional:    true,
				Computed:    true,
				Description: "List of domain policies in organization",
			},
			privacyPolicyVar: {
				Type:        schema.TypeSet,
				Elem:        GetPrivacyPolicyDatasource(),
				Optional:    true,
				Computed:    true,
				Description: "List of privacy policies in organization",
			},

			usersVar: {
				Type:        schema.TypeSet,
				Elem:        GetUserDatasource(),
				Optional:    true,
				Computed:    true,
				Description: "List of users in organization",
			},
			projectsVar: {
				Type:        schema.TypeSet,
				Elem:        GetProjectDatasource(),
				Optional:    true,
				Computed:    true,
				Description: "List of projects in organization",
			},
			domainsVar: {
				Type:        schema.TypeSet,
				Elem:        GetDomainDatasource(),
				Optional:    true,
				Computed:    true,
				Description: "List of domains in organization",
			},
			actionsVar: {
				Type:        schema.TypeSet,
				Elem:        GetActionDatasource(),
				Optional:    true,
				Computed:    true,
				Description: "List of actions in organization",
			},
		},
		ReadContext: readOrg,
	}
}

func readOrg(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started read")

	clientinfo, err := GetClientInfo(d)
	if err != nil {
		return diag.FromErr(err)
	}

	client, err := getManagementClient(clientinfo, "")
	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := client.GetMyOrg(ctx, &management2.GetMyOrgRequest{})
	if err != nil {
		return diag.Errorf("failed to get org: %v", err)
	}
	id := resp.GetOrg().GetId()
	d.SetId(id)
	name := resp.GetOrg().GetName()

	tflog.Debug(ctx, "found org", map[string]interface{}{
		"id":   id,
		"name": name,
	})

	if err := d.Set(nameVar, name); err != nil {
		return diag.Errorf("failed to set org name: %v", err)
	}
	if err := d.Set(orgVar, id); err != nil {
		return diag.Errorf("failed to set org: %v", err)
	}
	d.SetId(id)

	/****************************************************************************************
	Users
	*/
	users := d.Get(usersVar).(*schema.Set)
	if err := readUsersOfOrg(ctx, users, m, clientinfo, resp.GetOrg().GetId()); err != nil {
		return err
	}
	if err := d.Set(usersVar, users); err != nil {
		return diag.Errorf("failed to set list of users: %v", err)
	}

	/****************************************************************************************
	Projects
	*/
	projects := d.Get(projectsVar).(*schema.Set)
	if err := readProjectsOfOrg(ctx, projects, m, clientinfo, resp.GetOrg().GetId()); err != nil {
		return err
	}
	if err := d.Set(projectsVar, projects); err != nil {
		return diag.Errorf("failed to set list of projects: %v", err)
	}
	/****************************************************************************************
	Domains
	*/
	domains := d.Get(domainsVar).(*schema.Set)
	if err := readDomainsOfOrg(ctx, domains, m, clientinfo, resp.GetOrg().GetId()); err != nil {
		return err
	}
	if err := d.Set(domainsVar, domains); err != nil {
		return diag.Errorf("failed to set list of domains: %v", err)
	}

	/****************************************************************************************
	iam policy
	*/
	iamPolicy := d.Get(iamPolicyVar).(*schema.Set)
	if err := readIAMPolicyOfOrg(ctx, domains, m, clientinfo, resp.GetOrg().GetId()); err != nil {
		return err
	}
	if err := d.Set(iamPolicyVar, iamPolicy); err != nil {
		return diag.Errorf("failed to set list of iam policies: %v", err)
	}

	/****************************************************************************************
	label policy
	*/
	labelPolicy := d.Get(labelPolicyVar).(*schema.Set)
	if err := readLabelPolicyOfOrg(ctx, domains, m, clientinfo, resp.GetOrg().GetId()); err != nil {
		return err
	}
	if err := d.Set(labelPolicyVar, labelPolicy); err != nil {
		return diag.Errorf("failed to set list of label policies: %v", err)
	}

	/****************************************************************************************
	lockout policy
	*/
	lockoutPolicy := d.Get(lockoutPolicyVar).(*schema.Set)
	if err := readLockoutPolicyOfOrg(ctx, domains, m, clientinfo, resp.GetOrg().GetId()); err != nil {
		return err
	}
	if err := d.Set(lockoutPolicyVar, lockoutPolicy); err != nil {
		return diag.Errorf("failed to set list of lockout policies: %v", err)
	}

	/****************************************************************************************
	login policy
	*/
	loginPolicy := d.Get(loginPolicyVar).(*schema.Set)
	if err := readLoginPolicyOfOrg(ctx, domains, m, clientinfo, resp.GetOrg().GetId()); err != nil {
		return err
	}
	if err := d.Set(loginPolicyVar, loginPolicy); err != nil {
		return diag.Errorf("failed to set list of login policies: %v", err)
	}

	/****************************************************************************************
	password complexity policy
	*/
	passwordComplexityPolicy := d.Get(passwordComplexityPolicyVar).(*schema.Set)
	if err := readPasswordComplexityPolicyPolicyOfOrg(ctx, domains, m, clientinfo, resp.GetOrg().GetId()); err != nil {
		return err
	}
	if err := d.Set(passwordComplexityPolicyVar, passwordComplexityPolicy); err != nil {
		return diag.Errorf("failed to set list of password complexity policies: %v", err)
	}

	/****************************************************************************************
	privacy policy
	*/
	privacyPolicy := d.Get(privacyPolicyVar).(*schema.Set)
	if err := readPrivacyPolicyOfOrg(ctx, domains, m, clientinfo, resp.GetOrg().GetId()); err != nil {
		return err
	}
	if err := d.Set(privacyPolicyVar, privacyPolicy); err != nil {
		return diag.Errorf("failed to set list of privacy policies: %v", err)
	}

	/****************************************************************************************
	actions
	*/
	actions := d.Get(actionsVar).(*schema.Set)
	if err := readActionsOfOrg(ctx, domains, m, clientinfo, resp.GetOrg().GetId()); err != nil {
		return err
	}
	if err := d.Set(actionsVar, actions); err != nil {
		return diag.Errorf("failed to set list of actions: %v", err)
	}

	return nil
}

func resourceToValueMap(r *schema.Resource, d *schema.ResourceData) map[string]interface{} {
	values := make(map[string]interface{}, 0)
	for key := range r.Schema {
		values[key] = d.Get(key)
	}
	return values
}
