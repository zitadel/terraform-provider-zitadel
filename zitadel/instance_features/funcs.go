package instance_features

import (
	"context"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	feature "github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/feature/v2"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
)

func delete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started delete")

	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := helper.GetFeatureClient(ctx, clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.ResetInstanceFeatures(ctx, &feature.ResetInstanceFeaturesRequest{})
	if err != nil {
		return diag.Errorf("failed to reset instance features: %v", err)
	}
	return nil
}

func update(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started update")

	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := helper.GetFeatureClient(ctx, clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}

	req := &feature.SetInstanceFeaturesRequest{}

	if d.HasChange(loginDefaultOrgVar) {
		v := d.Get(loginDefaultOrgVar).(bool)
		req.LoginDefaultOrg = &v
	}

	if d.HasChange(userSchemaVar) {
		v := d.Get(userSchemaVar).(bool)
		req.UserSchema = &v
	}

	if d.HasChange(oidcTokenExchangeVar) {
		v := d.Get(oidcTokenExchangeVar).(bool)
		req.OidcTokenExchange = &v
	}

	if d.HasChange(improvedPerformanceVar) {
		set := d.Get(improvedPerformanceVar).(*schema.Set)
		req.ImprovedPerformance = make([]feature.ImprovedPerformance, 0, set.Len())
		for _, v := range set.List() {
			req.ImprovedPerformance = append(req.ImprovedPerformance, mapToImprovedPerformanceEnum(v.(string)))
		}
	}

	if d.HasChange(debugOidcParentErrorVar) {
		v := d.Get(debugOidcParentErrorVar).(bool)
		req.DebugOidcParentError = &v
	}

	if d.HasChange(oidcSingleV1SessionTerminationVar) {
		v := d.Get(oidcSingleV1SessionTerminationVar).(bool)
		req.OidcSingleV1SessionTermination = &v
	}

	if d.HasChange(enableBackChannelLogoutVar) {
		v := d.Get(enableBackChannelLogoutVar).(bool)
		req.EnableBackChannelLogout = &v
	}

	if d.HasChange(loginV2Var) {
		if v, ok := d.GetOk(loginV2Var); ok {
			list := v.([]interface{})
			if len(list) > 0 && list[0] != nil {
				loginV2Map := list[0].(map[string]interface{})
				req.LoginV2 = &feature.LoginV2{
					Required: loginV2Map[loginV2RequiredVar].(bool),
				}
				if baseURI, ok := loginV2Map[loginV2BaseURIVar].(string); ok && baseURI != "" {
					req.LoginV2.BaseUri = &baseURI
				}
			}
		}
	}

	if d.HasChange(permissionCheckV2Var) {
		v := d.Get(permissionCheckV2Var).(bool)
		req.PermissionCheckV2 = &v
	}

	if d.HasChange(consoleUseV2UserAPIVar) {
		v := d.Get(consoleUseV2UserAPIVar).(bool)
		req.ConsoleUseV2UserApi = &v
	}

	_, err = client.SetInstanceFeatures(ctx, req)
	if err != nil {
		return diag.Errorf("failed to update instance features: %v", err)
	}
	return nil
}

func create(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started create")

	d.SetId("instance")

	return update(ctx, d, m)
}

func read(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started read")

	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := helper.GetFeatureClient(ctx, clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := client.GetInstanceFeatures(ctx, &feature.GetInstanceFeaturesRequest{
		Inheritance: false,
	})

	if err != nil && helper.IgnoreIfNotFoundError(err) == nil {
		d.SetId("")
		return nil
	}
	if err != nil {
		return diag.Errorf("failed to get instance features: %v", err)
	}

	set := map[string]interface{}{}

	if resp.LoginDefaultOrg != nil {
		set[loginDefaultOrgVar] = resp.LoginDefaultOrg.Enabled
	}

	if resp.UserSchema != nil {
		set[userSchemaVar] = resp.UserSchema.Enabled
	}

	if resp.OidcTokenExchange != nil {
		set[oidcTokenExchangeVar] = resp.OidcTokenExchange.Enabled
	}

	if resp.ImprovedPerformance != nil {
		perfSet := schema.NewSet(schema.HashString, []interface{}{})
		for _, v := range resp.ImprovedPerformance.ExecutionPaths {
			perfSet.Add(mapFromImprovedPerformanceEnum(v))
		}
		set[improvedPerformanceVar] = perfSet
	}

	if resp.DebugOidcParentError != nil {
		set[debugOidcParentErrorVar] = resp.DebugOidcParentError.Enabled
	}

	if resp.OidcSingleV1SessionTermination != nil {
		set[oidcSingleV1SessionTerminationVar] = resp.OidcSingleV1SessionTermination.Enabled
	}

	if resp.EnableBackChannelLogout != nil {
		set[enableBackChannelLogoutVar] = resp.EnableBackChannelLogout.Enabled
	}

	if resp.LoginV2 != nil {
		loginV2Map := map[string]interface{}{
			loginV2RequiredVar: resp.LoginV2.Required,
		}
		if resp.LoginV2.BaseUri != nil {
			loginV2Map[loginV2BaseURIVar] = *resp.LoginV2.BaseUri
		}
		set[loginV2Var] = []interface{}{loginV2Map}
	}

	if resp.PermissionCheckV2 != nil {
		set[permissionCheckV2Var] = resp.PermissionCheckV2.Enabled
	}

	if resp.ConsoleUseV2UserApi != nil {
		set[consoleUseV2UserAPIVar] = resp.ConsoleUseV2UserApi.Enabled
	}

	for k, v := range set {
		if err := d.Set(k, v); err != nil {
			return diag.Errorf("failed to set %s of instance features: %v", k, err)
		}
	}

	d.SetId("instance")
	return nil
}

func mapToImprovedPerformanceEnum(s string) feature.ImprovedPerformance {
	switch s {
	case improvedPerformanceProjectGrant:
		return feature.ImprovedPerformance_IMPROVED_PERFORMANCE_PROJECT_GRANT
	case improvedPerformanceProject:
		return feature.ImprovedPerformance_IMPROVED_PERFORMANCE_PROJECT
	case improvedPerformanceUserGrant:
		return feature.ImprovedPerformance_IMPROVED_PERFORMANCE_USER_GRANT
	case improvedPerformanceOrgDomainVerified:
		return feature.ImprovedPerformance_IMPROVED_PERFORMANCE_ORG_DOMAIN_VERIFIED
	default:
		return feature.ImprovedPerformance_IMPROVED_PERFORMANCE_UNSPECIFIED
	}
}

func mapFromImprovedPerformanceEnum(e feature.ImprovedPerformance) string {
	switch e {
	case feature.ImprovedPerformance_IMPROVED_PERFORMANCE_PROJECT_GRANT:
		return improvedPerformanceProjectGrant
	case feature.ImprovedPerformance_IMPROVED_PERFORMANCE_PROJECT:
		return improvedPerformanceProject
	case feature.ImprovedPerformance_IMPROVED_PERFORMANCE_USER_GRANT:
		return improvedPerformanceUserGrant
	case feature.ImprovedPerformance_IMPROVED_PERFORMANCE_ORG_DOMAIN_VERIFIED:
		return improvedPerformanceOrgDomainVerified
	default:
		return ""
	}
}
