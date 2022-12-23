package login_policy

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/idp"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/management"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/policy"
	"google.golang.org/protobuf/types/known/durationpb"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/helper"
)

func delete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started create")

	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	org := d.Get(orgIDVar).(string)
	client, err := helper.GetManagementClient(clientinfo, org)
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.ResetLoginPolicyToDefault(ctx, &management.ResetLoginPolicyToDefaultRequest{})
	if err != nil {
		return diag.Errorf("failed to reset login policy: %v", err)
	}
	return nil
}

func update(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started update")

	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	org := d.Get(orgIDVar).(string)
	client, err := helper.GetManagementClient(clientinfo, org)
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChanges(
		allowUsernamePasswordVar,
		allowRegisterVar,
		allowExternalIDPVar,
		forceMFAVar,
		passwordlessTypeVar,
		hidePasswordResetVar,
		ignoreUnknownUsernamesVar,
		defaultRedirectURIVar,
		passwordCheckLifetimeVar,
		externalLoginCheckLifetimeVar,
		mfaInitSkipLifetimeVar,
		secondFactorCheckLifetimeVar,
		multiFactorCheckLifetimeVar,
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
		_, err = client.UpdateCustomLoginPolicy(ctx, &management.UpdateCustomLoginPolicyRequest{
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
	}

	if d.HasChange(secondFactorsVar) {
		o, n := d.GetChange(secondFactorsVar)
		addSecondFactor, deleteSecondFactors := helper.GetAddAndDelete(helper.SetToStringSlice(o.(*schema.Set)), helper.SetToStringSlice(n.(*schema.Set)))

		for _, factor := range addSecondFactor {
			if _, err := client.AddSecondFactorToLoginPolicy(ctx, &management.AddSecondFactorToLoginPolicyRequest{
				Type: policy.SecondFactorType(policy.SecondFactorType_value[factor]),
			}); err != nil {
				return diag.FromErr(err)
			}
		}
		for _, factor := range deleteSecondFactors {
			if _, err := client.RemoveSecondFactorFromLoginPolicy(ctx, &management.RemoveSecondFactorFromLoginPolicyRequest{
				Type: policy.SecondFactorType(policy.SecondFactorType_value[factor]),
			}); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	if d.HasChange(multiFactorsVar) {
		o, n := d.GetChange(multiFactorsVar)
		addMultiFactor, deleteMultiFactors := helper.GetAddAndDelete(helper.SetToStringSlice(o.(*schema.Set)), helper.SetToStringSlice(n.(*schema.Set)))

		for _, factor := range addMultiFactor {
			if _, err := client.AddMultiFactorToLoginPolicy(ctx, &management.AddMultiFactorToLoginPolicyRequest{
				Type: policy.MultiFactorType(policy.MultiFactorType_value[factor]),
			}); err != nil {
				return diag.FromErr(err)
			}
		}
		for _, factor := range deleteMultiFactors {
			if _, err := client.RemoveMultiFactorFromLoginPolicy(ctx, &management.RemoveMultiFactorFromLoginPolicyRequest{
				Type: policy.MultiFactorType(policy.MultiFactorType_value[factor]),
			}); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	if d.HasChange(idpsVar) {
		o, n := d.GetChange(idpsVar)
		addIdps, deleteIdps := helper.GetAddAndDelete(helper.SetToStringSlice(o.(*schema.Set)), helper.SetToStringSlice(n.(*schema.Set)))

		for _, addIdp := range addIdps {
			var ownertype idp.IDPOwnerType
			_, err := client.GetOrgIDPByID(ctx, &management.GetOrgIDPByIDRequest{Id: addIdp})
			if err != nil {
				ownertype = idp.IDPOwnerType_IDP_OWNER_TYPE_SYSTEM
			} else {
				ownertype = idp.IDPOwnerType_IDP_OWNER_TYPE_ORG
			}
			if _, err := client.AddIDPToLoginPolicy(ctx, &management.AddIDPToLoginPolicyRequest{IdpId: addIdp, OwnerType: ownertype}); err != nil {
				return diag.FromErr(err)
			}
		}
		for _, deleteIdp := range deleteIdps {
			if _, err := client.RemoveIDPFromLoginPolicy(ctx, &management.RemoveIDPFromLoginPolicyRequest{IdpId: deleteIdp}); err != nil {
				return diag.FromErr(err)
			}
		}
	}
	return nil
}

func create(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started create")

	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	org := d.Get(orgIDVar).(string)
	client, err := helper.GetManagementClient(clientinfo, org)
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

	secondFactors := make([]policy.SecondFactorType, 0)
	if secondFactorsSet, ok := d.GetOk(secondFactorsVar); ok {
		for _, factor := range secondFactorsSet.(*schema.Set).List() {
			secondFactors = append(secondFactors, policy.SecondFactorType(policy.SecondFactorType_value[factor.(string)]))
		}
	}
	multiFactors := make([]policy.MultiFactorType, 0)
	if multiFactorsSet, ok := d.GetOk(multiFactorsVar); ok {
		for _, factor := range multiFactorsSet.(*schema.Set).List() {
			multiFactors = append(multiFactors, policy.MultiFactorType(policy.MultiFactorType_value[factor.(string)]))
		}
	}

	_, err = client.AddCustomLoginPolicy(ctx, &management.AddCustomLoginPolicyRequest{
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
		SecondFactors:              secondFactors,
		MultiFactors:               multiFactors,
	})
	if err != nil {
		return diag.Errorf("failed to create login policy: %v", err)
	}
	d.SetId(org)

	idps := helper.GetOkSetToStringSlice(d, idpsVar)
	for _, addIdp := range idps {
		var ownertype idp.IDPOwnerType
		_, err := client.GetOrgIDPByID(ctx, &management.GetOrgIDPByIDRequest{Id: addIdp})
		if err != nil {
			ownertype = idp.IDPOwnerType_IDP_OWNER_TYPE_SYSTEM
		} else {
			ownertype = idp.IDPOwnerType_IDP_OWNER_TYPE_ORG
		}
		if _, err := client.AddIDPToLoginPolicy(ctx, &management.AddIDPToLoginPolicyRequest{IdpId: addIdp, OwnerType: ownertype}); err != nil {
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

	org := d.Get(orgIDVar).(string)
	client, err := helper.GetManagementClient(clientinfo, org)
	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := client.GetLoginPolicy(ctx, &management.GetLoginPolicyRequest{})
	if err != nil && helper.IgnoreIfNotFoundError(err) == nil {
		d.SetId("")
		return nil
	}
	if err != nil {
		return diag.Errorf("failed to get login policy")
	}

	policy := resp.Policy
	if policy.GetIsDefault() == true {
		d.SetId("")
		return nil
	}
	set := map[string]interface{}{
		orgIDVar:                      policy.GetDetails().GetResourceOwner(),
		allowUsernamePasswordVar:      policy.GetAllowUsernamePassword(),
		allowRegisterVar:              policy.GetAllowRegister(),
		allowExternalIDPVar:           policy.GetAllowExternalIdp(),
		forceMFAVar:                   policy.GetForceMfa(),
		passwordlessTypeVar:           policy.GetPasswordlessType().String(),
		hidePasswordResetVar:          policy.GetHidePasswordReset(),
		ignoreUnknownUsernamesVar:     policy.GetIgnoreUnknownUsernames(),
		defaultRedirectURIVar:         policy.GetDefaultRedirectUri(),
		passwordCheckLifetimeVar:      policy.GetPasswordCheckLifetime().AsDuration().String(),
		externalLoginCheckLifetimeVar: policy.GetExternalLoginCheckLifetime().AsDuration().String(),
		mfaInitSkipLifetimeVar:        policy.GetMfaInitSkipLifetime().AsDuration().String(),
		secondFactorCheckLifetimeVar:  policy.GetSecondFactorCheckLifetime().AsDuration().String(),
		multiFactorCheckLifetimeVar:   policy.GetMultiFactorCheckLifetime().AsDuration().String(),
	}

	respSecond, err := client.ListLoginPolicySecondFactors(ctx, &management.ListLoginPolicySecondFactorsRequest{})
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

	respMulti, err := client.ListLoginPolicyMultiFactors(ctx, &management.ListLoginPolicyMultiFactorsRequest{})
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

	respIDPs, err := client.ListLoginPolicyIDPs(ctx, &management.ListLoginPolicyIDPsRequest{})
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
	d.SetId(policy.GetDetails().GetResourceOwner())
	return nil
}
