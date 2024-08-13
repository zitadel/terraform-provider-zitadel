package idp_saml

import (
	"context"

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
	client, err := helper.GetAdminClient(ctx, clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}
	resp, err := client.AddSAMLProvider(ctx, &admin.AddSAMLProviderRequest{
		Name:              idp_utils.StringValue(d, idp_utils.NameVar),
		Binding:           idp.SAMLBinding(idp.SAMLBinding_value[idp_utils.StringValue(d, BindingVar)]),
		WithSignedRequest: idp_utils.BoolValue(d, WithSignedRequestVar),
		ProviderOptions:   idp_utils.ProviderOptionsValue(d),
		Metadata:          &admin.AddSAMLProviderRequest_MetadataXml{MetadataXml: []byte(idp_utils.StringValue(d, MetadataXMLVar))},
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
	client, err := helper.GetAdminClient(ctx, clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}
	_, err = client.UpdateSAMLProvider(ctx, &admin.UpdateSAMLProviderRequest{
		Id:                d.Id(),
		Name:              idp_utils.StringValue(d, idp_utils.NameVar),
		Binding:           idp.SAMLBinding(idp.SAMLBinding_value[idp_utils.StringValue(d, BindingVar)]),
		WithSignedRequest: idp_utils.BoolValue(d, WithSignedRequestVar),
		ProviderOptions:   idp_utils.ProviderOptionsValue(d),
		Metadata:          &admin.UpdateSAMLProviderRequest_MetadataXml{MetadataXml: []byte(idp_utils.StringValue(d, MetadataXMLVar))},
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
	client, err := helper.GetAdminClient(ctx, clientinfo)
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
	idp := resp.GetIdp()
	cfg := idp.GetConfig()
	specificCfg := cfg.GetSaml()
	generalCfg := cfg.GetOptions()
	set := map[string]interface{}{
		idp_utils.NameVar:              idp.GetName(),
		MetadataXMLVar:                 string(specificCfg.GetMetadataXml()),
		BindingVar:                     specificCfg.GetBinding().String(),
		WithSignedRequestVar:           specificCfg.GetWithSignedRequest(),
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
	d.SetId(idp.Id)
	return nil
}
