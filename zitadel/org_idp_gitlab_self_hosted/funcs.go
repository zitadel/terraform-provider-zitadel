package org_idp_gitlab_self_hosted

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/management"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/idp_gitlab_self_hosted"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/idp_utils"
)

func create(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}
	client, err := helper.GetManagementClient(ctx, clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}
	resp, err := client.AddGitLabSelfHostedProvider(helper.CtxWithOrgID(ctx, d), &management.AddGitLabSelfHostedProviderRequest{
		Name:            idp_utils.StringValue(d, idp_utils.NameVar),
		ClientId:        idp_utils.StringValue(d, idp_utils.ClientIDVar),
		ClientSecret:    idp_utils.StringValue(d, idp_utils.ClientSecretVar),
		Scopes:          idp_utils.ScopesValue(d),
		ProviderOptions: idp_utils.ProviderOptionsValue(d),
		Issuer:          idp_utils.StringValue(d, idp_gitlab_self_hosted.IssuerVar),
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
	client, err := helper.GetManagementClient(ctx, clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}
	_, err = client.UpdateGitLabSelfHostedProvider(helper.CtxWithOrgID(ctx, d), &management.UpdateGitLabSelfHostedProviderRequest{
		Id:              d.Id(),
		Name:            idp_utils.StringValue(d, idp_utils.NameVar),
		ClientId:        idp_utils.StringValue(d, idp_utils.ClientIDVar),
		ClientSecret:    idp_utils.StringValue(d, idp_utils.ClientSecretVar),
		Scopes:          idp_utils.ScopesValue(d),
		ProviderOptions: idp_utils.ProviderOptionsValue(d),
		Issuer:          idp_utils.StringValue(d, idp_gitlab_self_hosted.IssuerVar),
	})
	if err != nil {
		return diag.Errorf("failed to update idp: %v", err)
	}
	return nil
}

func read(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}
	client, err := helper.GetManagementClient(ctx, clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}
	resp, err := client.GetProviderByID(helper.CtxWithOrgID(ctx, d), &management.GetProviderByIDRequest{Id: helper.GetID(d, idp_utils.IdpIDVar)})
	if err != nil && helper.IgnoreIfNotFoundError(err) == nil {
		d.SetId("")
		return nil
	}
	if err != nil {
		return diag.Errorf("failed to get idp")
	}
	idp := resp.GetIdp()
	cfg := idp.GetConfig()
	specificCfg := cfg.GetGitlabSelfHosted()
	generalCfg := cfg.GetOptions()
	set := map[string]interface{}{
		helper.OrgIDVar:                  idp.GetDetails().GetResourceOwner(),
		idp_utils.NameVar:                idp.GetName(),
		idp_utils.ClientIDVar:            specificCfg.GetClientId(),
		idp_utils.ClientSecretVar:        idp_utils.StringValue(d, idp_utils.ClientSecretVar),
		idp_utils.ScopesVar:              specificCfg.GetScopes(),
		idp_utils.IsLinkingAllowedVar:    generalCfg.GetIsLinkingAllowed(),
		idp_utils.IsCreationAllowedVar:   generalCfg.GetIsCreationAllowed(),
		idp_utils.IsAutoCreationVar:      generalCfg.GetIsAutoCreation(),
		idp_utils.IsAutoUpdateVar:        generalCfg.GetIsAutoUpdate(),
		idp_utils.AutoLinkingVar:         idp_utils.AutoLinkingString(generalCfg.GetAutoLinking()),
		idp_gitlab_self_hosted.IssuerVar: specificCfg.GetIssuer(),
	}
	for k, v := range set {
		if err := d.Set(k, v); err != nil {
			return diag.Errorf("failed to set %s of oidc idp: %v", k, err)
		}
	}
	d.SetId(idp.Id)
	return nil
}
