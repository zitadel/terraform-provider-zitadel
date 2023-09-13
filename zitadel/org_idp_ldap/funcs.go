package org_idp_ldap

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/idp"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/management"
	"google.golang.org/protobuf/types/known/durationpb"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/helper"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/idp_ldap"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/idp_utils"
)

func create(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}
	client, err := helper.GetManagementClient(clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}
	timeout, err := time.ParseDuration(d.Get(idp_ldap.TimeoutVar).(string))
	if err != nil {
		return diag.FromErr(err)
	}
	resp, err := client.AddLDAPProvider(helper.CtxWithOrgID(ctx, d), &management.AddLDAPProviderRequest{
		Name:            idp_utils.StringValue(d, idp_utils.NameVar),
		ProviderOptions: idp_utils.ProviderOptionsValue(d),

		Servers:           idp_utils.InterfaceToStringSlice(d.Get(idp_ldap.ServersVar)),
		StartTls:          idp_utils.BoolValue(d, idp_ldap.StartTLSVar),
		BaseDn:            idp_utils.StringValue(d, idp_ldap.BaseDNVar),
		BindDn:            idp_utils.StringValue(d, idp_ldap.BindDNVar),
		BindPassword:      idp_utils.StringValue(d, idp_ldap.BindPasswordVar),
		UserBase:          idp_utils.StringValue(d, idp_ldap.UserBaseVar),
		UserObjectClasses: helper.GetOkSetToStringSlice(d, idp_ldap.UserObjectClassesVar),
		UserFilters:       helper.GetOkSetToStringSlice(d, idp_ldap.UserFiltersVar),
		Timeout:           durationpb.New(timeout),

		Attributes: &idp.LDAPAttributes{
			IdAttribute:                idp_utils.StringValue(d, idp_ldap.IdAttributeVar),
			FirstNameAttribute:         idp_utils.StringValue(d, idp_ldap.FirstNameAttributeVar),
			LastNameAttribute:          idp_utils.StringValue(d, idp_ldap.LastNameAttributeVar),
			DisplayNameAttribute:       idp_utils.StringValue(d, idp_ldap.DisplayNameAttributeVar),
			NickNameAttribute:          idp_utils.StringValue(d, idp_ldap.NickNameAttributeVar),
			PreferredUsernameAttribute: idp_utils.StringValue(d, idp_ldap.PreferredUsernameAttributeVar),
			EmailAttribute:             idp_utils.StringValue(d, idp_ldap.EmailAttributeVar),
			EmailVerifiedAttribute:     idp_utils.StringValue(d, idp_ldap.EmailVerifiedAttributeVar),
			PhoneAttribute:             idp_utils.StringValue(d, idp_ldap.PhoneAttributeVar),
			PhoneVerifiedAttribute:     idp_utils.StringValue(d, idp_ldap.PhoneVerifiedAttributeVar),
			PreferredLanguageAttribute: idp_utils.StringValue(d, idp_ldap.PreferredLanguageAttributeVar),
			AvatarUrlAttribute:         idp_utils.StringValue(d, idp_ldap.AvatarURLAttributeVar),
			ProfileAttribute:           idp_utils.StringValue(d, idp_ldap.ProfileAttributeVar),
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
	client, err := helper.GetManagementClient(clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}
	timeout, err := time.ParseDuration(d.Get(idp_ldap.TimeoutVar).(string))
	if err != nil {
		return diag.FromErr(err)
	}
	_, err = client.UpdateLDAPProvider(helper.CtxWithOrgID(ctx, d), &management.UpdateLDAPProviderRequest{
		Id:              d.Id(),
		Name:            idp_utils.StringValue(d, idp_utils.NameVar),
		ProviderOptions: idp_utils.ProviderOptionsValue(d),

		Servers:           idp_utils.InterfaceToStringSlice(d.Get(idp_ldap.ServersVar)),
		StartTls:          idp_utils.BoolValue(d, idp_ldap.StartTLSVar),
		BaseDn:            idp_utils.StringValue(d, idp_ldap.BaseDNVar),
		BindDn:            idp_utils.StringValue(d, idp_ldap.BindDNVar),
		BindPassword:      idp_utils.StringValue(d, idp_ldap.BindPasswordVar),
		UserBase:          idp_utils.StringValue(d, idp_ldap.UserBaseVar),
		UserObjectClasses: helper.GetOkSetToStringSlice(d, idp_ldap.UserObjectClassesVar),
		UserFilters:       helper.GetOkSetToStringSlice(d, idp_ldap.UserFiltersVar),
		Timeout:           durationpb.New(timeout),

		Attributes: &idp.LDAPAttributes{
			IdAttribute:                idp_utils.StringValue(d, idp_ldap.IdAttributeVar),
			FirstNameAttribute:         idp_utils.StringValue(d, idp_ldap.FirstNameAttributeVar),
			LastNameAttribute:          idp_utils.StringValue(d, idp_ldap.LastNameAttributeVar),
			DisplayNameAttribute:       idp_utils.StringValue(d, idp_ldap.DisplayNameAttributeVar),
			NickNameAttribute:          idp_utils.StringValue(d, idp_ldap.NickNameAttributeVar),
			PreferredUsernameAttribute: idp_utils.StringValue(d, idp_ldap.PreferredUsernameAttributeVar),
			EmailAttribute:             idp_utils.StringValue(d, idp_ldap.EmailAttributeVar),
			EmailVerifiedAttribute:     idp_utils.StringValue(d, idp_ldap.EmailVerifiedAttributeVar),
			PhoneAttribute:             idp_utils.StringValue(d, idp_ldap.PhoneAttributeVar),
			PhoneVerifiedAttribute:     idp_utils.StringValue(d, idp_ldap.PhoneVerifiedAttributeVar),
			PreferredLanguageAttribute: idp_utils.StringValue(d, idp_ldap.PreferredLanguageAttributeVar),
			AvatarUrlAttribute:         idp_utils.StringValue(d, idp_ldap.AvatarURLAttributeVar),
			ProfileAttribute:           idp_utils.StringValue(d, idp_ldap.ProfileAttributeVar),
		},
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
	client, err := helper.GetManagementClient(clientinfo)
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
	specificCfg := cfg.GetLdap()
	attributesCfg := specificCfg.GetAttributes()
	generalCfg := cfg.GetOptions()
	set := map[string]interface{}{
		helper.OrgIDVar:                idp.GetDetails().GetResourceOwner(),
		idp_utils.NameVar:              idp.GetName(),
		idp_utils.IsLinkingAllowedVar:  generalCfg.GetIsLinkingAllowed(),
		idp_utils.IsCreationAllowedVar: generalCfg.GetIsCreationAllowed(),
		idp_utils.IsAutoCreationVar:    generalCfg.GetIsAutoCreation(),
		idp_utils.IsAutoUpdateVar:      generalCfg.GetIsAutoUpdate(),

		idp_ldap.ServersVar:           specificCfg.GetServers(),
		idp_ldap.StartTLSVar:          specificCfg.GetStartTls(),
		idp_ldap.BaseDNVar:            specificCfg.GetBaseDn(),
		idp_ldap.BindDNVar:            specificCfg.GetBindDn(),
		idp_ldap.BindPasswordVar:      idp_utils.StringValue(d, idp_ldap.BindPasswordVar),
		idp_ldap.UserBaseVar:          specificCfg.GetUserBase(),
		idp_ldap.UserObjectClassesVar: specificCfg.GetUserObjectClasses(),
		idp_ldap.UserFiltersVar:       specificCfg.GetUserFilters(),
		idp_ldap.TimeoutVar:           specificCfg.GetTimeout().AsDuration().String(),
		idp_ldap.IdAttributeVar:       attributesCfg.GetIdAttribute(),

		idp_ldap.FirstNameAttributeVar:         attributesCfg.GetFirstNameAttribute(),
		idp_ldap.LastNameAttributeVar:          attributesCfg.GetLastNameAttribute(),
		idp_ldap.DisplayNameAttributeVar:       attributesCfg.GetDisplayNameAttribute(),
		idp_ldap.NickNameAttributeVar:          attributesCfg.GetNickNameAttribute(),
		idp_ldap.PreferredUsernameAttributeVar: attributesCfg.GetPreferredUsernameAttribute(),
		idp_ldap.EmailAttributeVar:             attributesCfg.GetEmailAttribute(),
		idp_ldap.EmailVerifiedAttributeVar:     attributesCfg.GetEmailVerifiedAttribute(),
		idp_ldap.PhoneAttributeVar:             attributesCfg.GetPhoneAttribute(),
		idp_ldap.PhoneVerifiedAttributeVar:     attributesCfg.GetPhoneVerifiedAttribute(),
		idp_ldap.PreferredLanguageAttributeVar: attributesCfg.GetPreferredLanguageAttribute(),
		idp_ldap.AvatarURLAttributeVar:         attributesCfg.GetAvatarUrlAttribute(),
		idp_ldap.ProfileAttributeVar:           attributesCfg.GetProfileAttribute(),
	}
	for k, v := range set {
		if err := d.Set(k, v); err != nil {
			return diag.Errorf("failed to set %s of oidc idp: %v", k, err)
		}
	}
	d.SetId(idp.Id)
	return nil
}
