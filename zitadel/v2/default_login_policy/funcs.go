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
		if helper.IgnorePreconditionError(err) != nil {
			return diag.Errorf("failed to update login policy: %v", err)
		}
		if resp != nil {
			d.SetId(resp.GetDetails().GetResourceOwner())
		} else {
			resp, err := client.GetLoginPolicy(ctx, &admin.GetLoginPolicyRequest{})
			if err != nil {
				return diag.Errorf("failed to update default login policy: %v", err)
			}
			d.SetId(resp.GetPolicy().GetDetails().GetResourceOwner())
		}
	}

	if d.HasChange(secondFactorsVar) {
		o, err := client.ListLoginPolicySecondFactors(ctx, &admin.ListLoginPolicySecondFactorsRequest{})
		if err != nil {
			return diag.Errorf("failed to get default login policy second factors: %v", err)
		}
		factors := make([]string, len(o.GetResult()))
		for i, factor := range o.GetResult() {
			factors[i] = policy.SecondFactorType_name[int32(factor.Number())]
		}
		addSecondFactor, deleteSecondFactors := helper.GetAddAndDelete(factors, helper.SetToStringSlice(d.Get(secondFactorsVar).(*schema.Set)))

		for _, factor := range addSecondFactor {
			if _, err := client.AddSecondFactorToLoginPolicy(ctx, &admin.AddSecondFactorToLoginPolicyRequest{
				Type: policy.SecondFactorType(policy.SecondFactorType_value[factor]),
			}); helper.IgnoreAlreadyExistsError(err) != nil {
				return diag.FromErr(err)
			}
		}
		for _, factor := range deleteSecondFactors {
			if _, err := client.RemoveSecondFactorFromLoginPolicy(ctx, &admin.RemoveSecondFactorFromLoginPolicyRequest{
				Type: policy.SecondFactorType(policy.SecondFactorType_value[factor]),
			}); helper.IgnoreAlreadyExistsError(err) != nil {
				return diag.FromErr(err)
			}
		}
	}

	if d.HasChange(multiFactorsVar) {
		o, err := client.ListLoginPolicyMultiFactors(ctx, &admin.ListLoginPolicyMultiFactorsRequest{})
		if err != nil {
			return diag.Errorf("failed to get default login policy multi factors: %v", err)
		}
		factors := make([]string, len(o.GetResult()))
		for i, factor := range o.GetResult() {
			factors[i] = policy.MultiFactorType_name[int32(factor.Number())]
		}
		addMultiFactor, deleteMultiFactors := helper.GetAddAndDelete(factors, helper.SetToStringSlice(d.Get(multiFactorsVar).(*schema.Set)))

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
		o, err := client.ListLoginPolicyIDPs(ctx, &admin.ListLoginPolicyIDPsRequest{})
		if err != nil {
			return diag.Errorf("failed to get default login policy idps: %v", err)
		}

		idps := make([]string, len(o.GetResult()))
		for i, idp := range o.GetResult() {
			idps[i] = idp.IdpId
		}
		addIdps, deleteIdps := helper.GetAddAndDelete(idps, helper.SetToStringSlice(d.Get(idpsVar).(*schema.Set)))

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
	if err != nil && helper.IgnoreIfNotFoundError(err) == nil {
		d.SetId("")
		return nil
	}
	if err != nil {
		return diag.Errorf("failed to get default login policy")
	}

	set := map[string]interface{}{
		allowUsernamePasswordVar:      resp.Policy.GetAllowUsernamePassword(),
		allowRegisterVar:              resp.Policy.GetAllowRegister(),
		allowExternalIDPVar:           resp.Policy.GetAllowExternalIdp(),
		forceMFAVar:                   resp.Policy.GetForceMfa(),
		passwordlessTypeVar:           resp.Policy.GetPasswordlessType().String(),
		hidePasswordResetVar:          resp.Policy.GetHidePasswordReset(),
		ignoreUnknownUsernamesVar:     resp.Policy.GetIgnoreUnknownUsernames(),
		defaultRedirectURIVar:         resp.Policy.GetDefaultRedirectUri(),
		passwordCheckLifetimeVar:      resp.Policy.GetPasswordCheckLifetime().AsDuration().String(),
		externalLoginCheckLifetimeVar: resp.Policy.GetExternalLoginCheckLifetime().AsDuration().String(),
		mfaInitSkipLifetimeVar:        resp.Policy.GetMfaInitSkipLifetime().AsDuration().String(),
		secondFactorCheckLifetimeVar:  resp.Policy.GetSecondFactorCheckLifetime().AsDuration().String(),
		multiFactorCheckLifetimeVar:   resp.Policy.GetMultiFactorCheckLifetime().AsDuration().String(),
	}

	respSecond, err := client.ListLoginPolicySecondFactors(ctx, &admin.ListLoginPolicySecondFactorsRequest{})
	if err != nil {
		return diag.Errorf("failed to get login policy secondfactors: %v", err)
	}
	if len(respSecond.GetResult()) > 0 {
		factors := make([]string, 0)
		for _, item := range respSecond.GetResult() {
			factors = append(factors, item.String())
		}
		set[secondFactorsVar] = factors
	}

	respMulti, err := client.ListLoginPolicyMultiFactors(ctx, &admin.ListLoginPolicyMultiFactorsRequest{})
	if err != nil {
		return diag.Errorf("failed to get login policy multifactors: %v", err)
	}
	if len(respMulti.GetResult()) > 0 {
		factors := make([]string, 0)
		for _, item := range respMulti.GetResult() {
			factors = append(factors, item.String())
		}
		set[multiFactorsVar] = factors
	}

	respIDPs, err := client.ListLoginPolicyIDPs(ctx, &admin.ListLoginPolicyIDPsRequest{})
	if err != nil {
		return diag.Errorf("failed to get login policy idps: %v", err)
	}
	if len(respIDPs.GetResult()) > 0 {
		idps := make([]string, 0)
		for _, idpItem := range respIDPs.GetResult() {
			idps = append(idps, idpItem.IdpId)
		}
		set[idpsVar] = idps
	}

	for k, v := range set {
		if err := d.Set(k, v); err != nil {
			return diag.Errorf("failed to set %s of login policy: %v", k, err)
		}
	}
	d.SetId(resp.Policy.GetDetails().GetResourceOwner())
	return nil
}
