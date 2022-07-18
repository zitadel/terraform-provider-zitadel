package v2

import (
	"context"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	management2 "github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/management"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/policy"
	"google.golang.org/protobuf/types/known/durationpb"
	"time"
)

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
	loginPolicyMultiFactorCheckLifetime   = "multi_factor_check_lifetime"
	loginPolicyIgnoreUnknownUsernames     = "ignore_unknown_usernames"
	loginPolicyDefaultRedirectURI         = "default_redirect_uri"
)

func GetLoginPolicy() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			loginPolicyOrgIdVar: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Id for the organization",
				ForceNew:    true,
			},
			loginPolicyAllowUsernamePassword: {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "defines if a user is allowed to login with his username and password",
			},
			loginPolicyAllowRegister: {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "defines if a person is allowed to register a user on this organisation",
			},
			loginPolicyAllowExternalIDP: {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "defines if a user is allowed to add a defined identity provider. E.g. Google auth",
			},
			loginPolicyForceMFA: {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "defines if a user MUST use a multi factor to log in",
			},
			loginPolicyPasswordlessType: {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "defines if passwordless is allowed for users",
			},
			loginPolicyIsDefault: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "defines if the organisation's admin changed the policy",
			},
			loginPolicyHidePasswordReset: {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "defines if password reset link should be shown in the login screen",
			},
			loginPolicyIgnoreUnknownUsernames: {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "defines if unknown username on login screen directly return an error or always display the password screen",
			},
			loginPolicyDefaultRedirectURI: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "defines where the user will be redirected to if the login is started without app context (e.g. from mail)",
			},
			loginPolicyPasswordCheckLifetime: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "",
			},
			loginPolicyExternalLoginCheckLifetime: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "",
			},
			loginPolicyMFAInitSkipLifetime: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "",
			},
			loginPolicySecondFactorCheckLifetime: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "",
			},
			loginPolicyMultiFactorCheckLifetime: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "",
			},
		},
		CreateContext: createLoginPolicy,
		UpdateContext: updateLoginPolicy,
		DeleteContext: deleteLoginPolicy,
		ReadContext:   readLoginPolicy,
	}
}

func deleteLoginPolicy(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started create")

	clientinfo, ok := m.(*ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	org := d.Get(loginPolicyOrgIdVar).(string)
	client, err := getManagementClient(clientinfo, org)
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.ResetLoginPolicyToDefault(ctx, &management2.ResetLoginPolicyToDefaultRequest{})
	if err != nil {
		return diag.Errorf("failed to reset login policy: %v", err)
	}
	d.SetId(org)
	return nil
}

func updateLoginPolicy(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started update")

	clientinfo, ok := m.(*ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	org := d.Get(loginPolicyOrgIdVar).(string)
	client, err := getManagementClient(clientinfo, org)
	if err != nil {
		return diag.FromErr(err)
	}

	passwordCheckLT, err := time.ParseDuration(d.Get(loginPolicyPasswordCheckLifetime).(string))
	if err != nil {
		return diag.FromErr(err)
	}
	externalLoginCheckLT, err := time.ParseDuration(d.Get(loginPolicyExternalLoginCheckLifetime).(string))
	if err != nil {
		return diag.FromErr(err)
	}
	mfaInitSkipLT, err := time.ParseDuration(d.Get(loginPolicyMFAInitSkipLifetime).(string))
	if err != nil {
		return diag.FromErr(err)
	}
	secondFactorCheckLT, err := time.ParseDuration(d.Get(loginPolicySecondFactorCheckLifetime).(string))
	if err != nil {
		return diag.FromErr(err)
	}
	multiFactorCheckLT, err := time.ParseDuration(d.Get(loginPolicyMultiFactorCheckLifetime).(string))
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.UpdateCustomLoginPolicy(ctx, &management2.UpdateCustomLoginPolicyRequest{
		AllowUsernamePassword:      d.Get(loginPolicyAllowUsernamePassword).(bool),
		AllowRegister:              d.Get(loginPolicyAllowRegister).(bool),
		AllowExternalIdp:           d.Get(loginPolicyAllowExternalIDP).(bool),
		ForceMfa:                   d.Get(loginPolicyForceMFA).(bool),
		PasswordlessType:           d.Get(loginPolicyPasswordlessType).(policy.PasswordlessType),
		HidePasswordReset:          d.Get(loginPolicyHidePasswordReset).(bool),
		IgnoreUnknownUsernames:     d.Get(loginPolicyIgnoreUnknownUsernames).(bool),
		DefaultRedirectUri:         d.Get(loginPolicyDefaultRedirectURI).(string),
		PasswordCheckLifetime:      durationpb.New(passwordCheckLT),
		ExternalLoginCheckLifetime: durationpb.New(externalLoginCheckLT),
		MfaInitSkipLifetime:        durationpb.New(mfaInitSkipLT),
		SecondFactorCheckLifetime:  durationpb.New(secondFactorCheckLT),
		MultiFactorCheckLifetime:   durationpb.New(multiFactorCheckLT),
	})
	if err != nil {
		return diag.Errorf("failed to update login policy: %v", err)
	}
	d.SetId(org)
	return nil
}

func createLoginPolicy(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started create")

	clientinfo, ok := m.(*ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	org := d.Get(loginPolicyOrgIdVar).(string)
	client, err := getManagementClient(clientinfo, org)
	if err != nil {
		return diag.FromErr(err)
	}

	passwordCheckLT, err := time.ParseDuration(d.Get(loginPolicyPasswordCheckLifetime).(string))
	if err != nil {
		return diag.FromErr(err)
	}
	externalLoginCheckLT, err := time.ParseDuration(d.Get(loginPolicyExternalLoginCheckLifetime).(string))
	if err != nil {
		return diag.FromErr(err)
	}
	mfaInitSkipLT, err := time.ParseDuration(d.Get(loginPolicyMFAInitSkipLifetime).(string))
	if err != nil {
		return diag.FromErr(err)
	}
	secondFactorCheckLT, err := time.ParseDuration(d.Get(loginPolicySecondFactorCheckLifetime).(string))
	if err != nil {
		return diag.FromErr(err)
	}
	multiFactorCheckLT, err := time.ParseDuration(d.Get(loginPolicyMultiFactorCheckLifetime).(string))
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.AddCustomLoginPolicy(ctx, &management2.AddCustomLoginPolicyRequest{
		AllowUsernamePassword:      d.Get(loginPolicyAllowUsernamePassword).(bool),
		AllowRegister:              d.Get(loginPolicyAllowRegister).(bool),
		AllowExternalIdp:           d.Get(loginPolicyAllowExternalIDP).(bool),
		ForceMfa:                   d.Get(loginPolicyForceMFA).(bool),
		PasswordlessType:           d.Get(loginPolicyPasswordlessType).(policy.PasswordlessType),
		HidePasswordReset:          d.Get(loginPolicyHidePasswordReset).(bool),
		IgnoreUnknownUsernames:     d.Get(loginPolicyIgnoreUnknownUsernames).(bool),
		DefaultRedirectUri:         d.Get(loginPolicyDefaultRedirectURI).(string),
		PasswordCheckLifetime:      durationpb.New(passwordCheckLT),
		ExternalLoginCheckLifetime: durationpb.New(externalLoginCheckLT),
		MfaInitSkipLifetime:        durationpb.New(mfaInitSkipLT),
		SecondFactorCheckLifetime:  durationpb.New(secondFactorCheckLT),
		MultiFactorCheckLifetime:   durationpb.New(multiFactorCheckLT),
	})
	if err != nil {
		return diag.Errorf("failed to create login policy: %v", err)
	}
	d.SetId(org)
	return nil
}

func readLoginPolicy(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started read")

	clientinfo, ok := m.(*ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	org := d.Get(domainPolicyOrgIdVar).(string)
	client, err := getManagementClient(clientinfo, org)
	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := client.GetLoginPolicy(ctx, &management2.GetLoginPolicyRequest{})
	if err != nil {
		return diag.Errorf("failed to get login policy: %v", err)
	}

	policy := resp.Policy
	set := map[string]interface{}{
		loginPolicyOrgIdVar:                   policy.GetDetails().GetResourceOwner(),
		loginPolicyIsDefault:                  policy.GetIsDefault(),
		loginPolicyAllowUsernamePassword:      policy.GetAllowUsernamePassword(),
		loginPolicyAllowRegister:              policy.GetAllowRegister(),
		loginPolicyAllowExternalIDP:           policy.GetAllowExternalIdp(),
		loginPolicyForceMFA:                   policy.GetForceMfa(),
		loginPolicyPasswordlessType:           policy.GetPasswordlessType(),
		loginPolicyHidePasswordReset:          policy.GetHidePasswordReset(),
		loginPolicyPasswordCheckLifetime:      policy.GetPasswordCheckLifetime(),
		loginPolicyExternalLoginCheckLifetime: policy.GetExternalLoginCheckLifetime(),
		loginPolicyMFAInitSkipLifetime:        policy.GetMfaInitSkipLifetime(),
		loginPolicySecondFactorCheckLifetime:  policy.GetSecondFactorCheckLifetime(),
		loginPolicyMultiFactorCheckLifetime:   policy.GetMultiFactorCheckLifetime(),
	}

	for k, v := range set {
		if err := d.Set(k, v); err != nil {
			return diag.Errorf("failed to set %s of login policy: %v", k, err)
		}
	}
	d.SetId(policy.GetDetails().GetResourceOwner())
	return nil
}
