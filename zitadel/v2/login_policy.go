package v2

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

const (
	loginPolicyOrgIdVar                   = "org_id"
	loginPolicyAllowUsernamePassword      = "user_login"
	loginPolicyAllowRegister              = "allow_register"
	loginPolicyAllowExternalIDP           = "allow_external_idp"
	loginPolicyForceMFA                   = "force_mfa"
	loginPolicyPasswordlessType           = "passwordless_type"
	loginPolicyIsDefault                  = "is_default"
	loginPolicyHidePasswordReset          = "hide_password_reset"
	loginPolicyPasswordCheckLifetime      = "password_check_lifetime"
	loginPolicyExternalLoginCheckLifetime = "external_login_check_lifetime"
	loginPolicyMFAInitSkipLifetime        = "mfa_init_skip_lifetime"
	loginPolicySecondFactorCheckLifetime  = "second_factor_check_lifetime"
	loginPolicyMultiFactorCHeckLifetime   = "multi_factor_check_lifetime"
)

func GetLoginPolicy() *schema.Resource {
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
			loginPolicyPasswordCheckLifetime: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "",
			},
			loginPolicyExternalLoginCheckLifetime: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "",
			},
			loginPolicyMFAInitSkipLifetime: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "",
			},
			loginPolicySecondFactorCheckLifetime: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "",
			},
			loginPolicyMultiFactorCHeckLifetime: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "",
			},
		},
	}
}
