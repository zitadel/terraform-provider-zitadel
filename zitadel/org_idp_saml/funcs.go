package org_idp_saml

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/idp"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/management"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/idp_saml"
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
	nameIdFormat := idp.SAMLNameIDFormat(idp.SAMLNameIDFormat_value[idp_utils.StringValue(d, idp_saml.NameIdFormatVar)])
	federatedLogoutEnabled := idp_utils.BoolValue(d, idp_saml.FederatedLogoutEnabledVar)
	req := &management.AddSAMLProviderRequest{
		Name:                          idp_utils.StringValue(d, idp_utils.NameVar),
		Binding:                       idp.SAMLBinding(idp.SAMLBinding_value[idp_utils.StringValue(d, idp_saml.BindingVar)]),
		WithSignedRequest:             idp_utils.BoolValue(d, idp_saml.WithSignedRequestVar),
		ProviderOptions:               idp_utils.ProviderOptionsValue(d),
		NameIdFormat:                  &nameIdFormat,
		TransientMappingAttributeName: helper.StringPtr(idp_utils.StringValue(d, idp_saml.TransientMappingAttributeNameVar)),
		FederatedLogoutEnabled:        &federatedLogoutEnabled,
		SignatureAlgorithm:            idp.SAMLSignatureAlgorithm(idp.SAMLSignatureAlgorithm_value[idp_utils.StringValue(d, idp_saml.SignatureAlgorithmVar)]),
	}
	if v, ok := d.GetOk(idp_saml.MetadataURLVar); ok && v.(string) != "" {
		req.Metadata = &management.AddSAMLProviderRequest_MetadataUrl{MetadataUrl: v.(string)}
	} else {
		req.Metadata = &management.AddSAMLProviderRequest_MetadataXml{MetadataXml: []byte(idp_utils.StringValue(d, idp_saml.MetadataXMLVar))}
	}
	resp, err := client.AddSAMLProvider(helper.CtxWithOrgID(ctx, d), req)
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
	nameIdFormat := idp.SAMLNameIDFormat(idp.SAMLNameIDFormat_value[idp_utils.StringValue(d, idp_saml.NameIdFormatVar)])
	federatedLogoutEnabled := idp_utils.BoolValue(d, idp_saml.FederatedLogoutEnabledVar)
	req := &management.UpdateSAMLProviderRequest{
		Id:                            d.Id(),
		Name:                          idp_utils.StringValue(d, idp_utils.NameVar),
		Binding:                       idp.SAMLBinding(idp.SAMLBinding_value[idp_utils.StringValue(d, idp_saml.BindingVar)]),
		WithSignedRequest:             idp_utils.BoolValue(d, idp_saml.WithSignedRequestVar),
		ProviderOptions:               idp_utils.ProviderOptionsValue(d),
		NameIdFormat:                  &nameIdFormat,
		TransientMappingAttributeName: helper.StringPtr(idp_utils.StringValue(d, idp_saml.TransientMappingAttributeNameVar)),
		FederatedLogoutEnabled:        &federatedLogoutEnabled,
		SignatureAlgorithm:            idp.SAMLSignatureAlgorithm(idp.SAMLSignatureAlgorithm_value[idp_utils.StringValue(d, idp_saml.SignatureAlgorithmVar)]),
	}
	if v, ok := d.GetOk(idp_saml.MetadataURLVar); ok && v.(string) != "" {
		req.Metadata = &management.UpdateSAMLProviderRequest_MetadataUrl{MetadataUrl: v.(string)}
	} else {
		req.Metadata = &management.UpdateSAMLProviderRequest_MetadataXml{MetadataXml: []byte(idp_utils.StringValue(d, idp_saml.MetadataXMLVar))}
	}
	_, err = client.UpdateSAMLProvider(helper.CtxWithOrgID(ctx, d), req)
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
		return diag.Errorf("failed to get idp: %v", err)
	}
	idp := resp.GetIdp()
	cfg := idp.GetConfig()
	specificCfg := cfg.GetSaml()
	generalCfg := cfg.GetOptions()
	set := map[string]interface{}{
		helper.OrgIDVar:                           idp.GetDetails().GetResourceOwner(),
		idp_utils.NameVar:                         idp.GetName(),
		idp_saml.BindingVar:                       specificCfg.GetBinding().String(),
		idp_saml.WithSignedRequestVar:             specificCfg.GetWithSignedRequest(),
		idp_saml.NameIdFormatVar:                  specificCfg.GetNameIdFormat().String(),
		idp_saml.TransientMappingAttributeNameVar: specificCfg.GetTransientMappingAttributeName(),
		idp_saml.FederatedLogoutEnabledVar:        specificCfg.GetFederatedLogoutEnabled(),
		idp_saml.SignatureAlgorithmVar:            specificCfg.GetSignatureAlgorithm().String(),
		idp_utils.IsLinkingAllowedVar:             generalCfg.GetIsLinkingAllowed(),
		idp_utils.IsCreationAllowedVar:            generalCfg.GetIsCreationAllowed(),
		idp_utils.IsAutoCreationVar:               generalCfg.GetIsAutoCreation(),
		idp_utils.IsAutoUpdateVar:                 generalCfg.GetIsAutoUpdate(),
		idp_utils.AutoLinkingVar:                  idp_utils.AutoLinkingString(generalCfg.GetAutoLinking()),
	}
	// Only set metadata_xml if the user did not configure metadata_url,
	// otherwise the resolved XML would cause a perpetual diff.
	if _, urlSet := d.GetOk(idp_saml.MetadataURLVar); !urlSet {
		set[idp_saml.MetadataXMLVar] = string(specificCfg.GetMetadataXml())
	}
	for k, v := range set {
		if err := d.Set(k, v); err != nil {
			return diag.Errorf("failed to set %s of idp: %v", k, err)
		}
	}
	d.SetId(idp.Id)
	return nil
}
