package org_idp_oidc

import (
	"context"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/idp"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/management"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/helper"
)

func delete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started delete")

	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := helper.GetManagementClient(clientinfo, d.Get(orgIDVar).(string))
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.RemoveOrgIDP(ctx, &management.RemoveOrgIDPRequest{
		IdpId: d.Id(),
	})
	if err != nil {
		return diag.Errorf("failed to delete oidc idp: %v", err)
	}
	return nil
}

func create(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started create")

	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := helper.GetManagementClient(clientinfo, d.Get(orgIDVar).(string))
	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := client.AddOrgOIDCIDP(ctx, &management.AddOrgOIDCIDPRequest{
		Name:               d.Get(nameVar).(string),
		StylingType:        idp.IDPStylingType(idp.IDPStylingType_value[d.Get(stylingTypeVar).(string)]),
		ClientId:           d.Get(clientIDVar).(string),
		ClientSecret:       d.Get(clientSecretVar).(string),
		Issuer:             d.Get(issuerVar).(string),
		Scopes:             helper.GetOkSetToStringSlice(d, scopesVar),
		DisplayNameMapping: idp.OIDCMappingField(idp.OIDCMappingField_value[d.Get(displayNameMappingVar).(string)]),
		UsernameMapping:    idp.OIDCMappingField(idp.OIDCMappingField_value[d.Get(usernameMappingVar).(string)]),
		AutoRegister:       d.Get(autoRegisterVar).(bool),
	})
	if err != nil {
		return diag.Errorf("failed to create oidc idp: %v", err)
	}
	d.SetId(resp.GetIdpId())

	return nil
}

func update(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started update")

	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := helper.GetManagementClient(clientinfo, d.Get(orgIDVar).(string))
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChanges(nameVar, stylingTypeVar, autoRegisterVar) {
		_, err := client.UpdateOrgIDP(ctx, &management.UpdateOrgIDPRequest{
			IdpId:        d.Id(),
			Name:         d.Get(nameVar).(string),
			StylingType:  idp.IDPStylingType(idp.IDPStylingType_value[d.Get(stylingTypeVar).(string)]),
			AutoRegister: d.Get(autoRegisterVar).(bool),
		})
		if err != nil {
			return diag.Errorf("failed to update oidc idp: %v", err)
		}
	}

	if d.HasChanges(clientIDVar, clientSecretVar, issuerVar, displayNameMappingVar, usernameMappingVar, scopesVar) {
		_, err = client.UpdateOrgIDPOIDCConfig(ctx, &management.UpdateOrgIDPOIDCConfigRequest{
			IdpId:              d.Id(),
			ClientId:           d.Get(clientIDVar).(string),
			ClientSecret:       d.Get(clientSecretVar).(string),
			Issuer:             d.Get(issuerVar).(string),
			Scopes:             helper.GetOkSetToStringSlice(d, scopesVar),
			DisplayNameMapping: idp.OIDCMappingField(idp.OIDCMappingField_value[d.Get(displayNameMappingVar).(string)]),
			UsernameMapping:    idp.OIDCMappingField(idp.OIDCMappingField_value[d.Get(usernameMappingVar).(string)]),
		})
		if err != nil {
			return diag.Errorf("failed to update oidc idp config: %v", err)
		}

	}
	return nil
}

func read(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, "started read")

	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}

	client, err := helper.GetManagementClient(clientinfo, d.Get(orgIDVar).(string))
	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := client.GetOrgIDPByID(ctx, &management.GetOrgIDPByIDRequest{Id: helper.GetID(d, idpIDVar)})
	if err != nil && helper.IgnoreIfNotFoundError(err) == nil {
		d.SetId("")
		return nil
	}
	if err != nil {
		return diag.Errorf("failed to get org idp oidc")
	}

	idp := resp.GetIdp()
	oidc := idp.GetOidcConfig()
	set := map[string]interface{}{
		orgIDVar:              idp.GetDetails().GetResourceOwner(),
		nameVar:               idp.GetName(),
		stylingTypeVar:        idp.GetStylingType().String(),
		clientIDVar:           oidc.GetClientId(),
		clientSecretVar:       d.Get(clientSecretVar).(string),
		issuerVar:             oidc.GetIssuer(),
		scopesVar:             oidc.GetScopes(),
		displayNameMappingVar: oidc.GetDisplayNameMapping().String(),
		usernameMappingVar:    oidc.GetUsernameMapping().String(),
		autoRegisterVar:       idp.GetAutoRegister(),
	}
	for k, v := range set {
		if err := d.Set(k, v); err != nil {
			return diag.Errorf("failed to set %s of oidc idp: %v", k, err)
		}
	}
	d.SetId(idp.Id)

	return nil
}
