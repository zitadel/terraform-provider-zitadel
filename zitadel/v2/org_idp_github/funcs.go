package org_idp_github

import (
	"context"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/org_idp_utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/idp"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/management"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/helper"
)

func create(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}
	client, err := helper.GetManagementClient(clientinfo, d.Get(org_idp_utils.OrgIDVar).(string))
	if err != nil {
		return diag.FromErr(err)
	}
	resp, err := client.AddGitHubProvider(ctx, &management.AddGitHubProviderRequest{
		Name:         d.Get(org_idp_utils.NameVar).(string),
		ClientId:     d.Get(org_idp_utils.ClientIDVar).(string),
		ClientSecret: d.Get(org_idp_utils.ClientSecretVar).(string),
		Scopes:       helper.GetOkSetToStringSlice(d, org_idp_utils.ScopesVar),
		ProviderOptions: &idp.Options{
			IsLinkingAllowed:  d.Get(org_idp_utils.IsLinkingAllowedVar).(bool),
			IsCreationAllowed: d.Get(org_idp_utils.IsCreationAllowedVar).(bool),
			IsAutoUpdate:      d.Get(org_idp_utils.IsAutoUpdateVar).(bool),
			IsAutoCreation:    d.Get(org_idp_utils.IsAutoCreationVar).(bool),
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
	client, err := helper.GetManagementClient(clientinfo, d.Get(org_idp_utils.OrgIDVar).(string))
	if err != nil {
		return diag.FromErr(err)
	}
	if d.HasChangesExcept(org_idp_utils.IdpIDVar, org_idp_utils.OrgIDVar) {
		_, err = client.UpdateGitHubProvider(ctx, &management.UpdateGitHubProviderRequest{
			Id:           d.Id(),
			Name:         d.Get(org_idp_utils.NameVar).(string),
			ClientId:     d.Get(org_idp_utils.ClientIDVar).(string),
			ClientSecret: d.Get(org_idp_utils.ClientSecretVar).(string),
			Scopes:       helper.GetOkSetToStringSlice(d, org_idp_utils.ScopesVar),
			ProviderOptions: &idp.Options{
				IsLinkingAllowed:  d.Get(org_idp_utils.IsLinkingAllowedVar).(bool),
				IsCreationAllowed: d.Get(org_idp_utils.IsCreationAllowedVar).(bool),
				IsAutoCreation:    d.Get(org_idp_utils.IsAutoCreationVar).(bool),
				IsAutoUpdate:      d.Get(org_idp_utils.IsAutoUpdateVar).(bool),
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
	client, err := helper.GetManagementClient(clientinfo, d.Get(org_idp_utils.OrgIDVar).(string))
	if err != nil {
		return diag.FromErr(err)
	}
	resp, err := client.GetProviderByID(ctx, &management.GetProviderByIDRequest{Id: helper.GetID(d, org_idp_utils.IdpIDVar)})
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
		org_idp_utils.OrgIDVar:             idp.GetDetails().GetResourceOwner(),
		org_idp_utils.NameVar:              idp.GetName(),
		org_idp_utils.ClientIDVar:          specificCfg.GetClientId(),
		org_idp_utils.ClientSecretVar:      d.Get(org_idp_utils.ClientSecretVar).(string),
		org_idp_utils.ScopesVar:            specificCfg.GetScopes(),
		org_idp_utils.IsLinkingAllowedVar:  generalCfg.GetIsLinkingAllowed(),
		org_idp_utils.IsCreationAllowedVar: generalCfg.GetIsCreationAllowed(),
		org_idp_utils.IsAutoCreationVar:    generalCfg.GetIsAutoCreation(),
		org_idp_utils.IsAutoUpdateVar:      generalCfg.GetIsAutoUpdate(),
	}
	for k, v := range set {
		if err := d.Set(k, v); err != nil {
			return diag.Errorf("failed to set %s of oidc idp: %v", k, err)
		}
	}
	d.SetId(idp.Id)
	return nil
}
