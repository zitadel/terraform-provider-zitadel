package org_idp_azure_ad

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/idp"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/management"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/helper"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/idp_azure_ad"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/idp_utils"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/org_idp_utils"
)

func create(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}
	client, err := helper.GetManagementClient(clientinfo, idp_utils.StringValue(d, org_idp_utils.OrgIDVar))
	if err != nil {
		return diag.FromErr(err)
	}
	resp, err := client.AddAzureADProvider(ctx, &management.AddAzureADProviderRequest{
		Name:            idp_utils.StringValue(d, idp_utils.NameVar),
		ClientId:        idp_utils.StringValue(d, idp_utils.ClientIDVar),
		ClientSecret:    idp_utils.StringValue(d, idp_utils.ClientSecretVar),
		Scopes:          idp_utils.ScopesValue(d),
		ProviderOptions: idp_utils.ProviderOptionsValue(d),
		Tenant:          idp_azure_ad.ConstructTenant(d),
		EmailVerified:   idp_utils.BoolValue(d, idp_azure_ad.EmailVerifiedVar),
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
	client, err := helper.GetManagementClient(clientinfo, idp_utils.StringValue(d, org_idp_utils.OrgIDVar))
	if err != nil {
		return diag.FromErr(err)
	}
	_, err = client.UpdateAzureADProvider(ctx, &management.UpdateAzureADProviderRequest{
		Id:              d.Id(),
		Name:            idp_utils.StringValue(d, idp_utils.NameVar),
		ClientId:        idp_utils.StringValue(d, idp_utils.ClientIDVar),
		ClientSecret:    idp_utils.StringValue(d, idp_utils.ClientSecretVar),
		Scopes:          idp_utils.ScopesValue(d),
		ProviderOptions: idp_utils.ProviderOptionsValue(d),
		Tenant:          idp_azure_ad.ConstructTenant(d),
		EmailVerified:   idp_utils.BoolValue(d, idp_azure_ad.EmailVerifiedVar),
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
	client, err := helper.GetManagementClient(clientinfo, idp_utils.StringValue(d, org_idp_utils.OrgIDVar))
	if err != nil {
		return diag.FromErr(err)
	}
	resp, err := client.GetProviderByID(ctx, &management.GetProviderByIDRequest{Id: helper.GetID(d, idp_utils.IdpIDVar)})
	if err != nil && helper.IgnoreIfNotFoundError(err) == nil {
		d.SetId("")
		return nil
	}
	if err != nil {
		return diag.Errorf("failed to get idp")
	}
	respIdp := resp.GetIdp()
	cfg := respIdp.GetConfig()
	specificCfg := cfg.GetAzureAd()
	generalCfg := cfg.GetOptions()
	set := map[string]interface{}{
		org_idp_utils.OrgIDVar:         respIdp.GetDetails().GetResourceOwner(),
		idp_utils.NameVar:              respIdp.GetName(),
		idp_utils.ClientIDVar:          specificCfg.GetClientId(),
		idp_utils.ClientSecretVar:      idp_utils.StringValue(d, idp_utils.ClientSecretVar),
		idp_utils.ScopesVar:            specificCfg.GetScopes(),
		idp_utils.IsLinkingAllowedVar:  generalCfg.GetIsLinkingAllowed(),
		idp_utils.IsCreationAllowedVar: generalCfg.GetIsCreationAllowed(),
		idp_utils.IsAutoCreationVar:    generalCfg.GetIsAutoCreation(),
		idp_utils.IsAutoUpdateVar:      generalCfg.GetIsAutoUpdate(),
		idp_azure_ad.EmailVerifiedVar:  specificCfg.GetEmailVerified(),
		idp_azure_ad.TenantTypeVar:     idp.AzureADTenantType_name[int32(specificCfg.GetTenant().GetTenantType())],
		idp_azure_ad.TenantIDVar:       specificCfg.GetTenant().GetTenantId(),
	}
	for k, v := range set {
		if err := d.Set(k, v); err != nil {
			return diag.Errorf("failed to set %s of oidc idp: %v", k, err)
		}
	}
	d.SetId(respIdp.Id)
	return nil
}
