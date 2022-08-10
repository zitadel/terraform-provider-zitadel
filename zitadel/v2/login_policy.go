package v2

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/idp"
	management2 "github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/management"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/policy"
	"google.golang.org/protobuf/types/known/durationpb"
)

const (
	loginPolicyOrgIdVar                   = "org_id"
	loginPolicyAllowUsernamePassword      = "user_login"
	loginPolicyAllowRegister              = "allow_register"
	loginPolicyAllowExternalIDP           = "allow_external_idp"
	loginPolicyForceMFA                   = "force_mfa"
	loginPolicyPasswordlessType           = "passwordless_type"
	loginPolicyHidePasswordReset          = "hide_password_reset"
	loginPolicyPasswordCheckLifetime      = "password_check_lifetime"
	loginPolicyExternalLoginCheckLifetime = "external_login_check_lifetime"
	loginPolicyMFAInitSkipLifetime        = "mfa_init_skip_lifetime"
	loginPolicySecondFactorCheckLifetime  = "second_factor_check_lifetime"
	loginPolicyMultiFactorCheckLifetime   = "multi_factor_check_lifetime"
	loginPolicyIgnoreUnknownUsernames     = "ignore_unknown_usernames"
	loginPolicyDefaultRedirectURI         = "default_redirect_uri"
	loginPolicySecondFactorsVar           = "second_factors"
	loginPolicyMultiFactorsVar            = "multi_factors"
	loginPolicyIDPsVar                    = "idps"
)

func GetLoginPolicy() *schema.Resource {
	return &schema.Resource{
		Description: "Resource representing the custom login policy of an organization.",
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
				Type:        schema.TypeString,
				Required:    true,
				Description: "defines if passwordless is allowed for users",
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
			loginPolicySecondFactorsVar: {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Required:    true,
				Description: "allowed second factors",
			},
			loginPolicyMultiFactorsVar: {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Required:    true,
				Description: "allowed multi factors",
			},
			loginPolicyIDPsVar: {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Required:    true,
				Description: "allowed idps to login or register",
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

	current, err := client.GetLoginPolicy(ctx, &management2.GetLoginPolicyRequest{})
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

	allowUsernamePassword := d.Get(loginPolicyAllowUsernamePassword).(bool)
	allowRegister := d.Get(loginPolicyAllowRegister).(bool)
	allowExternalIdp := d.Get(loginPolicyAllowExternalIDP).(bool)
	forceMfa := d.Get(loginPolicyForceMFA).(bool)
	passwordlessType := policy.PasswordlessType(policy.PasswordlessType_value[d.Get(loginPolicyPasswordlessType).(string)])
	hidePasswordReset := d.Get(loginPolicyHidePasswordReset).(bool)
	ignoreUnkownUsernames := d.Get(loginPolicyIgnoreUnknownUsernames).(bool)
	defaultRedirectUri := d.Get(loginPolicyDefaultRedirectURI).(string)
	currentPolicy := current.GetPolicy()
	if currentPolicy.GetAllowUsernamePassword() != allowUsernamePassword ||
		currentPolicy.GetAllowRegister() != allowRegister ||
		currentPolicy.GetAllowExternalIdp() != allowExternalIdp ||
		currentPolicy.GetForceMfa() != forceMfa ||
		currentPolicy.GetPasswordlessType() != passwordlessType ||
		currentPolicy.GetHidePasswordReset() != hidePasswordReset ||
		currentPolicy.GetIgnoreUnknownUsernames() != ignoreUnkownUsernames ||
		currentPolicy.GetDefaultRedirectUri() != defaultRedirectUri {

		_, err = client.UpdateCustomLoginPolicy(ctx, &management2.UpdateCustomLoginPolicyRequest{
			AllowUsernamePassword:      allowUsernamePassword,
			AllowRegister:              allowRegister,
			AllowExternalIdp:           allowExternalIdp,
			ForceMfa:                   forceMfa,
			PasswordlessType:           passwordlessType,
			HidePasswordReset:          hidePasswordReset,
			IgnoreUnknownUsernames:     ignoreUnkownUsernames,
			DefaultRedirectUri:         defaultRedirectUri,
			PasswordCheckLifetime:      durationpb.New(passwordCheckLT),
			ExternalLoginCheckLifetime: durationpb.New(externalLoginCheckLT),
			MfaInitSkipLifetime:        durationpb.New(mfaInitSkipLT),
			SecondFactorCheckLifetime:  durationpb.New(secondFactorCheckLT),
			MultiFactorCheckLifetime:   durationpb.New(multiFactorCheckLT),
		})
		if err != nil {
			return diag.Errorf("failed to update login policy: %v", err)
		}
	}
	d.SetId(org)

	secondFactors := setToStringSlice(d.Get(loginPolicySecondFactorsVar).(*schema.Set))
	currentSecondFactors := make([]stringify, 0)
	for _, secondFactor := range current.GetPolicy().GetSecondFactors() {
		currentSecondFactors = append(currentSecondFactors, secondFactor)
	}
	addSecondFactor, deleteSecondFactors := getAddAndDelete(currentSecondFactors, secondFactors)

	for _, factor := range addSecondFactor {
		if _, err := client.AddSecondFactorToLoginPolicy(ctx, &management2.AddSecondFactorToLoginPolicyRequest{
			Type: policy.SecondFactorType(policy.SecondFactorType_value[factor]),
		}); err != nil {
			return diag.FromErr(err)
		}
	}
	for _, factor := range deleteSecondFactors {
		if _, err := client.RemoveSecondFactorFromLoginPolicy(ctx, &management2.RemoveSecondFactorFromLoginPolicyRequest{
			Type: policy.SecondFactorType(policy.SecondFactorType_value[factor]),
		}); err != nil {
			return diag.FromErr(err)
		}
	}

	multiFactors := setToStringSlice(d.Get(loginPolicyMultiFactorsVar).(*schema.Set))
	currentMultiFactors := make([]stringify, 0)
	for _, multiFactor := range current.GetPolicy().GetMultiFactors() {
		currentMultiFactors = append(currentMultiFactors, multiFactor)
	}
	addMultiFactor, deleteMultiFactors := getAddAndDelete(currentMultiFactors, multiFactors)
	for _, factor := range addMultiFactor {
		if _, err := client.AddMultiFactorToLoginPolicy(ctx, &management2.AddMultiFactorToLoginPolicyRequest{
			Type: policy.MultiFactorType(policy.MultiFactorType_value[factor]),
		}); err != nil {
			return diag.FromErr(err)
		}
	}
	for _, factor := range deleteMultiFactors {
		if _, err := client.RemoveMultiFactorFromLoginPolicy(ctx, &management2.RemoveMultiFactorFromLoginPolicyRequest{
			Type: policy.MultiFactorType(policy.MultiFactorType_value[factor]),
		}); err != nil {
			return diag.FromErr(err)
		}
	}

	idps := setToStringSlice(d.Get(loginPolicyIDPsVar).(*schema.Set))
	currentIdps := make([]stringify, 0)
	for _, currentIdp := range current.GetPolicy().GetIdps() {
		currentIdps = append(currentIdps, &stringified{currentIdp.IdpId})
	}
	addIdps, deleteIdps := getAddAndDelete(currentIdps, idps)
	for _, addIdp := range addIdps {
		var ownertype idp.IDPOwnerType
		_, err := client.GetOrgIDPByID(ctx, &management2.GetOrgIDPByIDRequest{Id: addIdp})
		if err != nil {
			ownertype = idp.IDPOwnerType_IDP_OWNER_TYPE_SYSTEM
		} else {
			ownertype = idp.IDPOwnerType_IDP_OWNER_TYPE_ORG
		}
		if _, err := client.AddIDPToLoginPolicy(ctx, &management2.AddIDPToLoginPolicyRequest{IdpId: addIdp, OwnerType: ownertype}); err != nil {
			return diag.FromErr(err)
		}
	}
	for _, deleteIdp := range deleteIdps {
		if _, err := client.RemoveIDPFromLoginPolicy(ctx, &management2.RemoveIDPFromLoginPolicyRequest{IdpId: deleteIdp}); err != nil {
			return diag.FromErr(err)
		}
	}

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
	secondFactors := make([]policy.SecondFactorType, 0)
	secondFactorsSet := d.Get(loginPolicySecondFactorsVar).(*schema.Set)
	for _, factor := range secondFactorsSet.List() {
		secondFactors = append(secondFactors, policy.SecondFactorType(policy.SecondFactorType_value[factor.(string)]))
	}
	multiFactors := make([]policy.MultiFactorType, 0)
	multiFactorsSet := d.Get(loginPolicyMultiFactorsVar).(*schema.Set)
	for _, factor := range multiFactorsSet.List() {
		multiFactors = append(multiFactors, policy.MultiFactorType(policy.MultiFactorType_value[factor.(string)]))
	}

	_, err = client.AddCustomLoginPolicy(ctx, &management2.AddCustomLoginPolicyRequest{
		AllowUsernamePassword:      d.Get(loginPolicyAllowUsernamePassword).(bool),
		AllowRegister:              d.Get(loginPolicyAllowRegister).(bool),
		AllowExternalIdp:           d.Get(loginPolicyAllowExternalIDP).(bool),
		ForceMfa:                   d.Get(loginPolicyForceMFA).(bool),
		PasswordlessType:           policy.PasswordlessType(policy.PasswordlessType_value[d.Get(loginPolicyPasswordlessType).(string)]),
		HidePasswordReset:          d.Get(loginPolicyHidePasswordReset).(bool),
		IgnoreUnknownUsernames:     d.Get(loginPolicyIgnoreUnknownUsernames).(bool),
		DefaultRedirectUri:         d.Get(loginPolicyDefaultRedirectURI).(string),
		PasswordCheckLifetime:      durationpb.New(passwordCheckLT),
		ExternalLoginCheckLifetime: durationpb.New(externalLoginCheckLT),
		MfaInitSkipLifetime:        durationpb.New(mfaInitSkipLT),
		SecondFactorCheckLifetime:  durationpb.New(secondFactorCheckLT),
		MultiFactorCheckLifetime:   durationpb.New(multiFactorCheckLT),
		SecondFactors:              secondFactors,
		MultiFactors:               multiFactors,
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
		d.SetId("")
		return nil
		//return diag.Errorf("failed to get login policy: %v", err)
	}

	policy := resp.Policy
	if policy.GetIsDefault() == true {
		d.SetId("")
		return nil
	}
	set := map[string]interface{}{
		loginPolicyOrgIdVar:                   policy.GetDetails().GetResourceOwner(),
		loginPolicyAllowUsernamePassword:      policy.GetAllowUsernamePassword(),
		loginPolicyAllowRegister:              policy.GetAllowRegister(),
		loginPolicyAllowExternalIDP:           policy.GetAllowExternalIdp(),
		loginPolicyForceMFA:                   policy.GetForceMfa(),
		loginPolicyPasswordlessType:           policy.GetPasswordlessType().String(),
		loginPolicyHidePasswordReset:          policy.GetHidePasswordReset(),
		loginPolicyPasswordCheckLifetime:      policy.GetPasswordCheckLifetime().AsDuration().String(),
		loginPolicyExternalLoginCheckLifetime: policy.GetExternalLoginCheckLifetime().AsDuration().String(),
		loginPolicyMFAInitSkipLifetime:        policy.GetMfaInitSkipLifetime().AsDuration().String(),
		loginPolicySecondFactorCheckLifetime:  policy.GetSecondFactorCheckLifetime().AsDuration().String(),
		loginPolicyMultiFactorCheckLifetime:   policy.GetMultiFactorCheckLifetime().AsDuration().String(),
	}

	for k, v := range set {
		if err := d.Set(k, v); err != nil {
			return diag.Errorf("failed to set %s of login policy: %v", k, err)
		}
	}
	d.SetId(policy.GetDetails().GetResourceOwner())
	return nil
}
