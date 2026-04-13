package idp_saml

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/admin"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/idp"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/idp_utils"
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
	nameIdFormat := idp.SAMLNameIDFormat(idp.SAMLNameIDFormat_value[idp_utils.StringValue(d, NameIdFormatVar)])
	federatedLogoutEnabled := idp_utils.BoolValue(d, FederatedLogoutEnabledVar)
	req := &admin.AddSAMLProviderRequest{
		Name:                          idp_utils.StringValue(d, idp_utils.NameVar),
		Binding:                       idp.SAMLBinding(idp.SAMLBinding_value[idp_utils.StringValue(d, BindingVar)]),
		WithSignedRequest:             idp_utils.BoolValue(d, WithSignedRequestVar),
		ProviderOptions:               idp_utils.ProviderOptionsValue(d),
		NameIdFormat:                  &nameIdFormat,
		TransientMappingAttributeName: helper.StringPtr(idp_utils.StringValue(d, TransientMappingAttributeNameVar)),
		FederatedLogoutEnabled:        &federatedLogoutEnabled,
		SignatureAlgorithm:            idp.SAMLSignatureAlgorithm(idp.SAMLSignatureAlgorithm_value[idp_utils.StringValue(d, SignatureAlgorithmVar)]),
	}
	if v, ok := d.GetOk(MetadataURLVar); ok && v.(string) != "" {
		req.Metadata = &admin.AddSAMLProviderRequest_MetadataUrl{MetadataUrl: v.(string)}
	} else {
		req.Metadata = &admin.AddSAMLProviderRequest_MetadataXml{MetadataXml: []byte(idp_utils.StringValue(d, MetadataXMLVar))}
	}
	resp, err := client.AddSAMLProvider(ctx, req)
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
	nameIdFormat := idp.SAMLNameIDFormat(idp.SAMLNameIDFormat_value[idp_utils.StringValue(d, NameIdFormatVar)])
	federatedLogoutEnabled := idp_utils.BoolValue(d, FederatedLogoutEnabledVar)
	req := &admin.UpdateSAMLProviderRequest{
		Id:                            d.Id(),
		Name:                          idp_utils.StringValue(d, idp_utils.NameVar),
		Binding:                       idp.SAMLBinding(idp.SAMLBinding_value[idp_utils.StringValue(d, BindingVar)]),
		WithSignedRequest:             idp_utils.BoolValue(d, WithSignedRequestVar),
		ProviderOptions:               idp_utils.ProviderOptionsValue(d),
		NameIdFormat:                  &nameIdFormat,
		TransientMappingAttributeName: helper.StringPtr(idp_utils.StringValue(d, TransientMappingAttributeNameVar)),
		FederatedLogoutEnabled:        &federatedLogoutEnabled,
		SignatureAlgorithm:            idp.SAMLSignatureAlgorithm(idp.SAMLSignatureAlgorithm_value[idp_utils.StringValue(d, SignatureAlgorithmVar)]),
	}
	if v, ok := d.GetOk(MetadataURLVar); ok && v.(string) != "" {
		req.Metadata = &admin.UpdateSAMLProviderRequest_MetadataUrl{MetadataUrl: v.(string)}
	} else {
		req.Metadata = &admin.UpdateSAMLProviderRequest_MetadataXml{MetadataXml: []byte(idp_utils.StringValue(d, MetadataXMLVar))}
	}
	_, err = client.UpdateSAMLProvider(ctx, req)
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
		return diag.Errorf("failed to get idp: %v", err)
	}
	idp := resp.GetIdp()
	cfg := idp.GetConfig()
	specificCfg := cfg.GetSaml()
	generalCfg := cfg.GetOptions()
	set := map[string]interface{}{
		idp_utils.NameVar:                idp.GetName(),
		BindingVar:                       specificCfg.GetBinding().String(),
		WithSignedRequestVar:             specificCfg.GetWithSignedRequest(),
		NameIdFormatVar:                  specificCfg.GetNameIdFormat().String(),
		TransientMappingAttributeNameVar: specificCfg.GetTransientMappingAttributeName(),
		FederatedLogoutEnabledVar:        specificCfg.GetFederatedLogoutEnabled(),
		SignatureAlgorithmVar:            specificCfg.GetSignatureAlgorithm().String(),
		idp_utils.IsLinkingAllowedVar:    generalCfg.GetIsLinkingAllowed(),
		idp_utils.IsCreationAllowedVar:   generalCfg.GetIsCreationAllowed(),
		idp_utils.IsAutoCreationVar:      generalCfg.GetIsAutoCreation(),
		idp_utils.IsAutoUpdateVar:        generalCfg.GetIsAutoUpdate(),
		idp_utils.AutoLinkingVar:         idp_utils.AutoLinkingString(generalCfg.GetAutoLinking()),
	}
	// Only set metadata_xml if the user did not configure metadata_url,
	// otherwise the resolved XML would cause a perpetual diff.
	if _, urlSet := d.GetOk(MetadataURLVar); !urlSet {
		set[MetadataXMLVar] = string(specificCfg.GetMetadataXml())
	}
	for k, v := range set {
		if err := d.Set(k, v); err != nil {
			return diag.Errorf("failed to set %s of oidc idp: %v", k, err)
		}
	}
	d.SetId(idp.Id)
	return nil
}
