package org_idp_jwt

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

	resp, err := client.AddOrgJWTIDP(ctx, &management.AddOrgJWTIDPRequest{
		Name:         d.Get(nameVar).(string),
		StylingType:  idp.IDPStylingType(idp.IDPStylingType_value[d.Get(stylingTypeVar).(string)]),
		JwtEndpoint:  d.Get(jwtEndpointVar).(string),
		Issuer:       d.Get(issuerVar).(string),
		KeysEndpoint: d.Get(keysEndpointVar).(string),
		HeaderName:   d.Get(headerNameVar).(string),
		AutoRegister: d.Get(autoRegisterVar).(bool),
	})
	if err != nil {
		return diag.Errorf("failed to create jwt idp: %v", err)
	}
	d.SetId(resp.IdpId)
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
			return diag.Errorf("failed to update jwt idp: %v", err)
		}
	}

	if d.HasChanges(jwtEndpointVar, issuerVar, keysEndpointVar, headerNameVar) {
		_, err = client.UpdateOrgIDPJWTConfig(ctx, &management.UpdateOrgIDPJWTConfigRequest{
			IdpId:        d.Id(),
			JwtEndpoint:  d.Get(jwtEndpointVar).(string),
			Issuer:       d.Get(issuerVar).(string),
			KeysEndpoint: d.Get(keysEndpointVar).(string),
			HeaderName:   d.Get(headerNameVar).(string),
		})
		if err != nil {
			return diag.Errorf("failed to update jwt idp config: %v", err)
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
		return diag.Errorf("failed to get org idp jwt")
	}

	idp := resp.GetIdp()
	jwt := idp.GetJwtConfig()
	set := map[string]interface{}{
		orgIDVar:        idp.GetDetails().ResourceOwner,
		nameVar:         idp.GetName(),
		stylingTypeVar:  idp.GetStylingType().String(),
		jwtEndpointVar:  jwt.GetJwtEndpoint(),
		issuerVar:       jwt.GetIssuer(),
		keysEndpointVar: jwt.GetKeysEndpoint(),
		headerNameVar:   jwt.GetHeaderName(),
		autoRegisterVar: idp.GetAutoRegister(),
	}
	for k, v := range set {
		if err := d.Set(k, v); err != nil {
			return diag.Errorf("failed to set %s of jwt idp: %v", k, err)
		}
	}
	d.SetId(idp.Id)
	return nil
}
