package idp_azure_ad

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/admin"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/idp"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/helper"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/idp_utils"
)

func create(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}
	client, err := helper.GetAdminClient(clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}
	resp, err := client.AddAzureADProvider(ctx, &admin.AddAzureADProviderRequest{
		Name:          d.Get(idp_utils.NameVar).(string),
		ClientId:      d.Get(idp_utils.ClientIDVar).(string),
		ClientSecret:  d.Get(idp_utils.ClientSecretVar).(string),
		Tenant:        constructTenant(d),
		EmailVerified: d.Get(idp_utils.EmailVerifiedVar).(bool),
		Scopes:        helper.GetOkSetToStringSlice(d, idp_utils.ScopesVar),
		ProviderOptions: &idp.Options{
			IsLinkingAllowed:  d.Get(idp_utils.IsLinkingAllowedVar).(bool),
			IsCreationAllowed: d.Get(idp_utils.IsCreationAllowedVar).(bool),
			IsAutoUpdate:      d.Get(idp_utils.IsAutoUpdateVar).(bool),
			IsAutoCreation:    d.Get(idp_utils.IsAutoCreationVar).(bool),
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
	client, err := helper.GetAdminClient(clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}
	if d.HasChangesExcept(idp_utils.IdpIDVar) {
		_, err = client.UpdateAzureADProvider(ctx, &admin.UpdateAzureADProviderRequest{
			Id:            d.Id(),
			Name:          d.Get(idp_utils.NameVar).(string),
			ClientId:      d.Get(idp_utils.ClientIDVar).(string),
			ClientSecret:  d.Get(idp_utils.ClientSecretVar).(string),
			Scopes:        helper.GetOkSetToStringSlice(d, idp_utils.ScopesVar),
			Tenant:        constructTenant(d),
			EmailVerified: d.Get(idp_utils.EmailVerifiedVar).(bool),
			ProviderOptions: &idp.Options{
				IsLinkingAllowed:  d.Get(idp_utils.IsLinkingAllowedVar).(bool),
				IsCreationAllowed: d.Get(idp_utils.IsCreationAllowedVar).(bool),
				IsAutoCreation:    d.Get(idp_utils.IsAutoCreationVar).(bool),
				IsAutoUpdate:      d.Get(idp_utils.IsAutoUpdateVar).(bool),
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
	client, err := helper.GetAdminClient(clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}
	resp, err := client.GetProviderByID(ctx, &admin.GetProviderByIDRequest{Id: helper.GetID(d, idp_utils.IdpIDVar)})
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
		idp_utils.NameVar:              respIdp.GetName(),
		idp_utils.ClientIDVar:          specificCfg.GetClientId(),
		idp_utils.ClientSecretVar:      d.Get(idp_utils.ClientSecretVar).(string),
		idp_utils.ScopesVar:            specificCfg.GetScopes(),
		idp_utils.EmailVerifiedVar:     specificCfg.GetEmailVerified(),
		idp_utils.TenantTypeVar:        idp.AzureADTenantType_name[int32(specificCfg.GetTenant().GetTenantType())],
		idp_utils.TenantIDVar:          specificCfg.GetTenant().GetTenantId(),
		idp_utils.IsLinkingAllowedVar:  generalCfg.GetIsLinkingAllowed(),
		idp_utils.IsCreationAllowedVar: generalCfg.GetIsCreationAllowed(),
		idp_utils.IsAutoCreationVar:    generalCfg.GetIsAutoCreation(),
		idp_utils.IsAutoUpdateVar:      generalCfg.GetIsAutoUpdate(),
	}
	for k, v := range set {
		if err := d.Set(k, v); err != nil {
			return diag.Errorf("failed to set %s of oidc idp: %v", k, err)
		}
	}
	d.SetId(respIdp.Id)
	return nil
}

func constructTenant(d *schema.ResourceData) *idp.AzureADTenant {
	tenant := &idp.AzureADTenant{}
	tenantId := d.Get(idp_utils.TenantIDVar).(string)
	if tenantId != "" {
		tenant.Type = &idp.AzureADTenant_TenantId{
			TenantId: tenantId,
		}
	} else {
		tenant.Type = &idp.AzureADTenant_TenantType{
			TenantType: idp.AzureADTenantType(idp.AzureADTenantType_value[d.Get(idp_utils.TenantTypeVar).(string)]),
		}
	}
	return tenant
}
