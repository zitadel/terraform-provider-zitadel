package default_login_policy

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/admin"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/policy"
	"google.golang.org/protobuf/types/known/durationpb"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/helper"
)

func delete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "default login policy cannot be deleted")
	return nil
}

func update(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started update")

	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := helper.GetAdminClient(clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}

	current, err := client.GetLoginPolicy(ctx, &admin.GetLoginPolicyRequest{})
	if err != nil {
		return diag.FromErr(err)
	}

	passwordCheckLT, err := time.ParseDuration(d.Get(passwordCheckLifetimeVar).(string))
	if err != nil {
		return diag.FromErr(err)
	}
	externalLoginCheckLT, err := time.ParseDuration(d.Get(externalLoginCheckLifetimeVar).(string))
	if err != nil {
		return diag.FromErr(err)
	}
	mfaInitSkipLT, err := time.ParseDuration(d.Get(mfaInitSkipLifetimeVar).(string))
	if err != nil {
		return diag.FromErr(err)
	}
	secondFactorCheckLT, err := time.ParseDuration(d.Get(secondFactorCheckLifetimeVar).(string))
	if err != nil {
		return diag.FromErr(err)
	}

	multiFactorCheckLT, err := time.ParseDuration(d.Get(multiFactorCheckLifetimeVar).(string))
	if err != nil {
		return diag.FromErr(err)
	}

	allowUsernamePassword := d.Get(allowUsernamePasswordVar).(bool)
	allowRegister := d.Get(allowRegisterVar).(bool)
	allowExternalIdp := d.Get(allowExternalIDPVar).(bool)
	forceMfa := d.Get(forceMFAVar).(bool)
	passwordlessType := policy.PasswordlessType(policy.PasswordlessType_value[d.Get(passwordlessTypeVar).(string)])
	hidePasswordReset := d.Get(hidePasswordResetVar).(bool)
	ignoreUnkownUsernames := d.Get(ignoreUnknownUsernamesVar).(bool)
	defaultRedirectUri := d.Get(defaultRedirectURIVar).(string)
	currentPolicy := current.GetPolicy()
	if currentPolicy.GetAllowUsernamePassword() != allowUsernamePassword ||
		currentPolicy.GetAllowRegister() != allowRegister ||
		currentPolicy.GetAllowExternalIdp() != allowExternalIdp ||
		currentPolicy.GetForceMfa() != forceMfa ||
		currentPolicy.GetPasswordlessType() != passwordlessType ||
		currentPolicy.GetHidePasswordReset() != hidePasswordReset ||
		currentPolicy.GetIgnoreUnknownUsernames() != ignoreUnkownUsernames ||
		currentPolicy.GetDefaultRedirectUri() != defaultRedirectUri {

		resp, err := client.UpdateLoginPolicy(ctx, &admin.UpdateLoginPolicyRequest{
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
		d.SetId(resp.GetDetails().GetResourceOwner())
	}

	secondFactors := helper.SetToStringSlice(d.Get(secondFactorsVar).(*schema.Set))
	currentSecondFactors := make([]helper.Stringify, 0)
	for _, secondFactor := range current.GetPolicy().GetSecondFactors() {
		currentSecondFactors = append(currentSecondFactors, secondFactor)
	}
	addSecondFactor, deleteSecondFactors := helper.GetAddAndDelete(currentSecondFactors, secondFactors)

	for _, factor := range addSecondFactor {
		if _, err := client.AddSecondFactorToLoginPolicy(ctx, &admin.AddSecondFactorToLoginPolicyRequest{
			Type: policy.SecondFactorType(policy.SecondFactorType_value[factor]),
		}); err != nil {
			return diag.FromErr(err)
		}
	}
	for _, factor := range deleteSecondFactors {
		if _, err := client.RemoveSecondFactorFromLoginPolicy(ctx, &admin.RemoveSecondFactorFromLoginPolicyRequest{
			Type: policy.SecondFactorType(policy.SecondFactorType_value[factor]),
		}); err != nil {
			return diag.FromErr(err)
		}
	}

	multiFactors := helper.SetToStringSlice(d.Get(multiFactorsVar).(*schema.Set))
	currentMultiFactors := make([]helper.Stringify, 0)
	for _, multiFactor := range current.GetPolicy().GetMultiFactors() {
		currentMultiFactors = append(currentMultiFactors, multiFactor)
	}
	addMultiFactor, deleteMultiFactors := helper.GetAddAndDelete(currentMultiFactors, multiFactors)
	for _, factor := range addMultiFactor {
		if _, err := client.AddMultiFactorToLoginPolicy(ctx, &admin.AddMultiFactorToLoginPolicyRequest{
			Type: policy.MultiFactorType(policy.MultiFactorType_value[factor]),
		}); err != nil {
			return diag.FromErr(err)
		}
	}
	for _, factor := range deleteMultiFactors {
		if _, err := client.RemoveMultiFactorFromLoginPolicy(ctx, &admin.RemoveMultiFactorFromLoginPolicyRequest{
			Type: policy.MultiFactorType(policy.MultiFactorType_value[factor]),
		}); err != nil {
			return diag.FromErr(err)
		}
	}

	idps := helper.SetToStringSlice(d.Get(idpsVar).(*schema.Set))
	currentIdps := make([]helper.Stringify, 0)
	for _, currentIdp := range current.GetPolicy().GetIdps() {
		currentIdps = append(currentIdps, &helper.Stringified{currentIdp.IdpId})
	}
	addIdps, deleteIdps := helper.GetAddAndDelete(currentIdps, idps)
	for _, addIdp := range addIdps {
		if _, err := client.AddIDPToLoginPolicy(ctx, &admin.AddIDPToLoginPolicyRequest{IdpId: addIdp}); err != nil {
			return diag.FromErr(err)
		}
	}
	for _, deleteIdp := range deleteIdps {
		if _, err := client.RemoveIDPFromLoginPolicy(ctx, &admin.RemoveIDPFromLoginPolicyRequest{IdpId: deleteIdp}); err != nil {
			return diag.FromErr(err)
		}
	}

	return nil
}

func read(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started read")

	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := helper.GetAdminClient(clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := client.GetLoginPolicy(ctx, &admin.GetLoginPolicyRequest{})
	if err != nil {
		d.SetId("")
		return nil
		//return diag.Errorf("failed to get login policy: %v", err)
	}

	policy := resp.Policy
	set := map[string]interface{}{
		allowUsernamePasswordVar:      policy.GetAllowUsernamePassword(),
		allowRegisterVar:              policy.GetAllowRegister(),
		allowExternalIDPVar:           policy.GetAllowExternalIdp(),
		forceMFAVar:                   policy.GetForceMfa(),
		passwordlessTypeVar:           policy.GetPasswordlessType().String(),
		hidePasswordResetVar:          policy.GetHidePasswordReset(),
		passwordCheckLifetimeVar:      policy.GetPasswordCheckLifetime().AsDuration().String(),
		externalLoginCheckLifetimeVar: policy.GetExternalLoginCheckLifetime().AsDuration().String(),
		mfaInitSkipLifetimeVar:        policy.GetMfaInitSkipLifetime().AsDuration().String(),
		secondFactorCheckLifetimeVar:  policy.GetSecondFactorCheckLifetime().AsDuration().String(),
		multiFactorCheckLifetimeVar:   policy.GetMultiFactorCheckLifetime().AsDuration().String(),
	}

	for k, v := range set {
		if err := d.Set(k, v); err != nil {
			return diag.Errorf("failed to set %s of login policy: %v", k, err)
		}
	}
	d.SetId(policy.GetDetails().GetResourceOwner())
	return nil
}
