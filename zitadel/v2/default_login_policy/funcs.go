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

	if d.HasChanges(passwordCheckLifetimeVar,
		externalLoginCheckLifetimeVar,
		mfaInitSkipLifetimeVar,
		secondFactorCheckLifetimeVar,
		multiFactorCheckLifetimeVar,
		allowUsernamePasswordVar,
		allowRegisterVar,
		allowExternalIDPVar,
		forceMFAVar,
		passwordlessTypeVar,
		hidePasswordResetVar,
		ignoreUnknownUsernamesVar,
		defaultRedirectURIVar,
	) {
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

		resp, err := client.UpdateLoginPolicy(ctx, &admin.UpdateLoginPolicyRequest{
			AllowUsernamePassword:      d.Get(allowUsernamePasswordVar).(bool),
			AllowRegister:              d.Get(allowRegisterVar).(bool),
			AllowExternalIdp:           d.Get(allowExternalIDPVar).(bool),
			ForceMfa:                   d.Get(forceMFAVar).(bool),
			PasswordlessType:           policy.PasswordlessType(policy.PasswordlessType_value[d.Get(passwordlessTypeVar).(string)]),
			HidePasswordReset:          d.Get(hidePasswordResetVar).(bool),
			IgnoreUnknownUsernames:     d.Get(ignoreUnknownUsernamesVar).(bool),
			DefaultRedirectUri:         d.Get(defaultRedirectURIVar).(string),
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

	if d.HasChange(secondFactorsVar) {
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
	}

	if d.HasChange(multiFactorsVar) {
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
	}

	if d.HasChange(idpsVar) {
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
	}

	set := map[string]interface{}{
		allowUsernamePasswordVar:      resp.Policy.GetAllowUsernamePassword(),
		allowRegisterVar:              resp.Policy.GetAllowRegister(),
		allowExternalIDPVar:           resp.Policy.GetAllowExternalIdp(),
		forceMFAVar:                   resp.Policy.GetForceMfa(),
		passwordlessTypeVar:           resp.Policy.GetPasswordlessType().String(),
		hidePasswordResetVar:          resp.Policy.GetHidePasswordReset(),
		passwordCheckLifetimeVar:      resp.Policy.GetPasswordCheckLifetime().AsDuration().String(),
		externalLoginCheckLifetimeVar: resp.Policy.GetExternalLoginCheckLifetime().AsDuration().String(),
		mfaInitSkipLifetimeVar:        resp.Policy.GetMfaInitSkipLifetime().AsDuration().String(),
		secondFactorCheckLifetimeVar:  resp.Policy.GetSecondFactorCheckLifetime().AsDuration().String(),
		multiFactorCheckLifetimeVar:   resp.Policy.GetMultiFactorCheckLifetime().AsDuration().String(),
	}

	secondFactors := &schema.Set{}
	for _, factor := range resp.Policy.SecondFactors {
		secondFactors.Add(policy.SecondFactorType_name[int32(factor.Number())])
	}
	set[secondFactorsVar] = secondFactors
	multiFactors := &schema.Set{}
	for _, factor := range resp.Policy.MultiFactors {
		multiFactors.Add(policy.MultiFactorType_name[int32(factor.Number())])
	}
	set[multiFactorsVar] = multiFactors
	idps := &schema.Set{}
	for _, idp := range resp.Policy.Idps {
		idps.Add(idp.IdpId)
	}
	set[idpsVar] = idps

	for k, v := range set {
		if err := d.Set(k, v); err != nil {
			return diag.Errorf("failed to set %s of login policy: %v", k, err)
		}
	}
	d.SetId(resp.Policy.GetDetails().GetResourceOwner())
	return nil
}
