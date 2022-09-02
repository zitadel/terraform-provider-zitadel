package idp_oidc

import (
	"context"
	"reflect"

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

	scopes := make([]string, 0)
	scopesSet := d.Get(scopesVar).(*schema.Set)
	for _, scope := range scopesSet.List() {
		scopes = append(scopes, scope.(string))
	}

	stylingType := d.Get(stylingTypeVar)
	displayNameMapping := d.Get(displayNameMappingVar).(string)
	usernameMapping := d.Get(usernameMappingVar).(string)
	resp, err := client.AddOrgOIDCIDP(ctx, &management.AddOrgOIDCIDPRequest{
		Name:               d.Get(nameVar).(string),
		StylingType:        idp.IDPStylingType(idp.IDPStylingType_value[stylingType.(string)]),
		ClientId:           d.Get(clientIDVar).(string),
		ClientSecret:       d.Get(clientSecretVar).(string),
		Issuer:             d.Get(issuerVar).(string),
		Scopes:             scopes,
		DisplayNameMapping: idp.OIDCMappingField(idp.OIDCMappingField_value[displayNameMapping]),
		UsernameMapping:    idp.OIDCMappingField(idp.OIDCMappingField_value[usernameMapping]),
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

	resp, err := client.GetOrgIDPByID(ctx, &management.GetOrgIDPByIDRequest{Id: d.Id()})
	if err != nil {
		return diag.Errorf("failed to read oidc idp: %v", err)
	}

	idpID := d.Id()
	name := d.Get(nameVar).(string)
	stylingType := d.Get(stylingTypeVar).(string)
	autoRegister := d.Get(autoRegisterVar).(bool)
	changed := false
	if resp.GetIdp().GetName() != name ||
		resp.GetIdp().GetStylingType().String() != stylingType ||
		resp.GetIdp().GetAutoRegister() != autoRegister {
		changed = true
		_, err := client.UpdateOrgIDP(ctx, &management.UpdateOrgIDPRequest{
			IdpId:        idpID,
			Name:         name,
			StylingType:  idp.IDPStylingType(idp.IDPStylingType_value[stylingType]),
			AutoRegister: autoRegister,
		})
		if err != nil {
			return diag.Errorf("failed to update oidc idp: %v", err)
		}
	}

	oidc := resp.GetIdp().GetOidcConfig()
	clientID := d.Get(clientIDVar).(string)
	clientSecret := d.Get(clientSecretVar).(string)
	issuer := d.Get(issuerVar).(string)
	scopesSet := d.Get(scopesVar).(*schema.Set)
	displayNameMapping := d.Get(displayNameMappingVar).(string)
	usernameMapping := d.Get(usernameMappingVar).(string)

	scopes := make([]string, 0)
	for _, scope := range scopesSet.List() {
		scopes = append(scopes, scope.(string))
	}

	//either nothing changed on the IDP or something besides the secret changed
	if (oidc.GetClientId() != clientID ||
		oidc.GetIssuer() != issuer ||
		!reflect.DeepEqual(oidc.GetScopes(), scopes) ||
		oidc.GetDisplayNameMapping().String() != displayNameMapping ||
		oidc.GetUsernameMapping().String() != usernameMapping) ||
		!changed {

		_, err = client.UpdateOrgIDPOIDCConfig(ctx, &management.UpdateOrgIDPOIDCConfigRequest{
			IdpId:              idpID,
			ClientId:           clientID,
			ClientSecret:       clientSecret,
			Issuer:             issuer,
			Scopes:             scopes,
			DisplayNameMapping: idp.OIDCMappingField(idp.OIDCMappingField_value[displayNameMapping]),
			UsernameMapping:    idp.OIDCMappingField(idp.OIDCMappingField_value[usernameMapping]),
		})
		if err != nil {
			return diag.Errorf("failed to update oidc idp config: %v", err)
		}
	}
	d.SetId(idpID)
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
	if err != nil {
		d.SetId("")
		return nil
		//return diag.Errorf("failed to read oidc idp: %v", err)
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
