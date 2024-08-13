package idp_azure_ad

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/admin"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/idp"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/helper"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/idp_utils"
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
	tenant, err := ConstructTenant(d)
	if err != nil {
		return diag.FromErr(err)
	}
	resp, err := client.AddAzureADProvider(ctx, &admin.AddAzureADProviderRequest{
		Name:            idp_utils.StringValue(d, idp_utils.NameVar),
		ClientId:        idp_utils.StringValue(d, idp_utils.ClientIDVar),
		ClientSecret:    idp_utils.StringValue(d, idp_utils.ClientSecretVar),
		Scopes:          idp_utils.ScopesValue(d),
		ProviderOptions: idp_utils.ProviderOptionsValue(d),
		Tenant:          tenant,
		EmailVerified:   idp_utils.BoolValue(d, EmailVerifiedVar),
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
	tenant, err := ConstructTenant(d)
	if err != nil {
		return diag.FromErr(err)
	}
	_, err = client.UpdateAzureADProvider(ctx, &admin.UpdateAzureADProviderRequest{
		Id:              d.Id(),
		Name:            idp_utils.StringValue(d, idp_utils.NameVar),
		ClientId:        idp_utils.StringValue(d, idp_utils.ClientIDVar),
		ClientSecret:    idp_utils.StringValue(d, idp_utils.ClientSecretVar),
		Scopes:          idp_utils.ScopesValue(d),
		ProviderOptions: idp_utils.ProviderOptionsValue(d),
		Tenant:          tenant,
		EmailVerified:   idp_utils.BoolValue(d, EmailVerifiedVar),
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
	tenantID := specificCfg.GetTenant().GetTenantId()
	set := map[string]interface{}{
		idp_utils.NameVar:              respIdp.GetName(),
		idp_utils.ClientIDVar:          specificCfg.GetClientId(),
		idp_utils.ClientSecretVar:      idp_utils.StringValue(d, idp_utils.ClientSecretVar),
		idp_utils.ScopesVar:            specificCfg.GetScopes(),
		idp_utils.IsLinkingAllowedVar:  generalCfg.GetIsLinkingAllowed(),
		idp_utils.IsCreationAllowedVar: generalCfg.GetIsCreationAllowed(),
		idp_utils.IsAutoCreationVar:    generalCfg.GetIsAutoCreation(),
		idp_utils.IsAutoUpdateVar:      generalCfg.GetIsAutoUpdate(),
		EmailVerifiedVar:               specificCfg.GetEmailVerified(),
		TenantIDVar:                    tenantID,
	}

	if tenantID == "" {
		set[TenantTypeVar] = idp.AzureADTenantType_name[int32(specificCfg.GetTenant().GetTenantType())]
	} else {
		set[TenantTypeVar] = ""
	}

	for k, v := range set {
		if err := d.Set(k, v); err != nil {
			return diag.Errorf("failed to set %s of oidc idp: %v", k, err)
		}
	}
	d.SetId(respIdp.Id)
	return nil
}

func ConstructTenant(d *schema.ResourceData) (*idp.AzureADTenant, error) {
	tenant := &idp.AzureADTenant{}
	tenantId := idp_utils.StringValue(d, TenantIDVar)
	tenantType := idp_utils.StringValue(d, TenantTypeVar)
	if tenantId == "" && tenantType == "" {
		return nil, fmt.Errorf("tenant_id or tenant_type are required, but both were empty")
	}
	if tenantId != "" && tenantType != "" {
		return nil, fmt.Errorf("tenant_id and tenant_type are mutually exclusive, but got id %s and type %s", tenantId, tenantType)
	}
	if tenantId != "" {
		tenant.Type = &idp.AzureADTenant_TenantId{
			TenantId: tenantId,
		}
	} else {
		tenant.Type = &idp.AzureADTenant_TenantType{
			TenantType: idp.AzureADTenantType(idp.AzureADTenantType_value[tenantType]),
		}
	}
	return tenant, nil
}
