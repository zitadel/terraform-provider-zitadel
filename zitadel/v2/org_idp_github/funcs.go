package org_idp_github

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/idp"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/management"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/helper"
)

func delete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}
	client, err := helper.GetManagementClient(clientinfo, d.Get(orgIDVar).(string))
	if err != nil {
		return diag.FromErr(err)
	}
	_, err = client.DeleteProvider(ctx, &management.DeleteProviderRequest{
		Id: d.Id(),
	})
	if err != nil {
		return diag.Errorf("failed to delete idp: %v", err)
	}
	return nil
}

func create(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}
	client, err := helper.GetManagementClient(clientinfo, d.Get(orgIDVar).(string))
	if err != nil {
		return diag.FromErr(err)
	}
	resp, err := client.AddGitHubProvider(ctx, &management.AddGitHubProviderRequest{
		Name:         d.Get(nameVar).(string),
		ClientId:     d.Get(clientIDVar).(string),
		ClientSecret: d.Get(clientSecretVar).(string),
		Scopes:       helper.GetOkSetToStringSlice(d, scopesVar),
		ProviderOptions: &idp.Options{
			IsLinkingAllowed:  d.Get(isLinkingAllowedVar).(bool),
			IsCreationAllowed: d.Get(isCreationAllowedVar).(bool),
			IsAutoUpdate:      d.Get(isAutoUpdateVar).(bool),
			IsAutoCreation:    d.Get(isAutoCreationVar).(bool),
		},
	})
	if err != nil {
		return diag.Errorf("failed to create idp: %v", err)
	}
	d.SetId(resp.GetId())
	return nil
}

func update(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}
	client, err := helper.GetManagementClient(clientinfo, d.Get(orgIDVar).(string))
	if err != nil {
		return diag.FromErr(err)
	}
	if d.HasChangesExcept(idpIDVar, orgIDVar) {
		_, err = client.UpdateGitHubProvider(ctx, &management.UpdateGitHubProviderRequest{
			Id:           d.Id(),
			Name:         d.Get(nameVar).(string),
			ClientId:     d.Get(clientIDVar).(string),
			ClientSecret: d.Get(clientSecretVar).(string),
			Scopes:       helper.GetOkSetToStringSlice(d, scopesVar),
			ProviderOptions: &idp.Options{
				IsLinkingAllowed:  d.Get(isLinkingAllowedVar).(bool),
				IsCreationAllowed: d.Get(isCreationAllowedVar).(bool),
				IsAutoCreation:    d.Get(isAutoCreationVar).(bool),
				IsAutoUpdate:      d.Get(isAutoUpdateVar).(bool),
			},
		})
		if err != nil {
			return diag.Errorf("failed to update idp: %v", err)
		}
	}
	return nil
}

func read(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}
	client, err := helper.GetManagementClient(clientinfo, d.Get(orgIDVar).(string))
	if err != nil {
		return diag.FromErr(err)
	}
	resp, err := client.GetProviderByID(ctx, &management.GetProviderByIDRequest{Id: helper.GetID(d, idpIDVar)})
	if err != nil && helper.IgnoreIfNotFoundError(err) == nil {
		d.SetId("")
		return nil
	}
	if err != nil {
		return diag.Errorf("failed to get idp")
	}
	idp := resp.GetIdp()
	cfg := idp.GetConfig()
	specificCfg := cfg.GetGithub()
	generalCfg := cfg.GetOptions()
	set := map[string]interface{}{
		orgIDVar:             idp.GetDetails().GetResourceOwner(),
		nameVar:              idp.GetName(),
		clientIDVar:          specificCfg.GetClientId(),
		clientSecretVar:      d.Get(clientSecretVar).(string),
		scopesVar:            specificCfg.GetScopes(),
		isLinkingAllowedVar:  generalCfg.GetIsLinkingAllowed(),
		isCreationAllowedVar: generalCfg.GetIsCreationAllowed(),
		isAutoCreationVar:    generalCfg.GetIsAutoCreation(),
		isAutoUpdateVar:      generalCfg.GetIsAutoUpdate(),
	}
	for k, v := range set {
		if err := d.Set(k, v); err != nil {
			return diag.Errorf("failed to set %s of oidc idp: %v", k, err)
		}
	}
	d.SetId(idp.Id)
	return nil
}
