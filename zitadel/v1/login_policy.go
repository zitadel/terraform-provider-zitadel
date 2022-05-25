package v1

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	management2 "github.com/zitadel/zitadel-go/pkg/client/zitadel/management"
)

const (
	loginPolicyOrgIdVar               = "org_id"
	loginPolicyAllowUsernamePassword  = "allow_username_password"
	loginPolicyAllowRegister          = "allow_register"
	loginPolicyAllowExternalIDP       = "allow_external_idp"
	loginPolicyForceMFA               = "force_mfa"
	loginPolicyPasswordlessType       = "passwordless_type"
	loginPolicyIsDefault              = "is_default"
	loginPolicyHidePasswordReset      = "hide_password_reset"
	loginPolicyIgnoreUnknownUsernames = "ignore_unknown_usernames"
	loginPolicyDefaultRedirectURI     = "default_redirect_uri"
)

func GetLoginPolicyDatasource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			loginPolicyOrgIdVar: {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Id for the organization",
			},
			loginPolicyAllowUsernamePassword: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "defines if a user is allowed to login with his username and password",
			},
			loginPolicyAllowRegister: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "defines if a person is allowed to register a user on this organisation",
			},
			loginPolicyAllowExternalIDP: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "defines if a user is allowed to add a defined identity provider. E.g. Google auth",
			},
			loginPolicyForceMFA: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "defines if a user MUST use a multi factor to log in",
			},
			loginPolicyPasswordlessType: {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "defines if passwordless is allowed for users",
			},
			loginPolicyIsDefault: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "defines if the organisation's admin changed the policy",
			},
			loginPolicyHidePasswordReset: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "defines if password reset link should be shown in the login screen",
			},
			loginPolicyIgnoreUnknownUsernames: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "defines if unknown username on login screen directly return an error or always display the password screen",
			},
			loginPolicyDefaultRedirectURI: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "defines where the user will be redirected to if the login is started without app context (e.g. from mail)",
			},
		},
	}
}

func readLoginPolicyOfOrg(ctx context.Context, policies *schema.Set, m interface{}, clientinfo *ClientInfo, org string) diag.Diagnostics {
	client, err := getManagementClient(clientinfo, org)
	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := client.GetLoginPolicy(ctx, &management2.GetLoginPolicyRequest{})
	if err != nil {
		return diag.Errorf("failed to get list of domains: %v", err)
	}

	policy := resp.Policy
	values := map[string]interface{}{
		loginPolicyOrgIdVar:               policy.GetDetails().GetResourceOwner(),
		loginPolicyAllowUsernamePassword:  policy.GetAllowUsernamePassword(),
		loginPolicyAllowRegister:          policy.GetAllowRegister(),
		loginPolicyAllowExternalIDP:       policy.GetAllowExternalIdp(),
		loginPolicyForceMFA:               policy.GetForceMfa(),
		loginPolicyPasswordlessType:       policy.GetPasswordlessType(),
		loginPolicyIsDefault:              policy.GetIsDefault(),
		loginPolicyHidePasswordReset:      policy.GetHidePasswordReset(),
		loginPolicyIgnoreUnknownUsernames: policy.GetIgnoreUnknownUsernames(),
		//loginPolicyDefaultRedirectURI: policy
	}
	policies.Add(values)
	return nil
}
