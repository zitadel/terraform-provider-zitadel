package org_idp_ldap

import (
	"context"
	"time"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/idp_ldap"

	"google.golang.org/protobuf/types/known/durationpb"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/idp_utils"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/org_idp_utils"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/idp"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/management"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/helper"
)

func create(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	clientinfo, ok := m.(*helper.ClientInfo)
	if !ok {
		return diag.Errorf("failed to get client")
	}
	client, err := helper.GetManagementClient(clientinfo, d.Get(org_idp_utils.OrgIDVar).(string))
	if err != nil {
		return diag.FromErr(err)
	}
	timeout, err := time.ParseDuration(d.Get(idp_ldap.TimeoutVar).(string))
	if err != nil {
		return diag.FromErr(err)
	}
	resp, err := client.AddLDAPProvider(ctx, &management.AddLDAPProviderRequest{
		Name: d.Get(idp_utils.NameVar).(string),
		ProviderOptions: &idp.Options{
			IsLinkingAllowed:  d.Get(idp_utils.IsLinkingAllowedVar).(bool),
			IsCreationAllowed: d.Get(idp_utils.IsCreationAllowedVar).(bool),
			IsAutoUpdate:      d.Get(idp_utils.IsAutoUpdateVar).(bool),
			IsAutoCreation:    d.Get(idp_utils.IsAutoCreationVar).(bool),
		},

		Servers:           idp_utils.InterfaceToStringSlice(d.Get(idp_ldap.ServersVar)),
		StartTls:          d.Get(idp_ldap.StartTLSVar).(bool),
		BaseDn:            d.Get(idp_ldap.BaseDNVar).(string),
		BindDn:            d.Get(idp_ldap.BindDNVar).(string),
		BindPassword:      d.Get(idp_ldap.BindPasswordVar).(string),
		UserBase:          d.Get(idp_ldap.UserBaseVar).(string),
		UserObjectClasses: helper.GetOkSetToStringSlice(d, idp_ldap.UserObjectClassesVar),
		UserFilters:       helper.GetOkSetToStringSlice(d, idp_ldap.UserFiltersVar),
		Timeout:           durationpb.New(timeout),

		Attributes: &idp.LDAPAttributes{
			IdAttribute:                d.Get(idp_ldap.IdAttributeVar).(string),
			FirstNameAttribute:         d.Get(idp_ldap.FirstNameAttributeVar).(string),
			LastNameAttribute:          d.Get(idp_ldap.LastNameAttributeVar).(string),
			DisplayNameAttribute:       d.Get(idp_ldap.DisplayNameAttributeVar).(string),
			NickNameAttribute:          d.Get(idp_ldap.NickNameAttributeVar).(string),
			PreferredUsernameAttribute: d.Get(idp_ldap.PreferredUsernameAttributeVar).(string),
			EmailAttribute:             d.Get(idp_ldap.EmailAttributeVar).(string),
			EmailVerifiedAttribute:     d.Get(idp_ldap.EmailVerifiedAttributeVar).(string),
			PhoneAttribute:             d.Get(idp_ldap.PhoneAttributeVar).(string),
			PhoneVerifiedAttribute:     d.Get(idp_ldap.PhoneVerifiedAttributeVar).(string),
			PreferredLanguageAttribute: d.Get(idp_ldap.PreferredLanguageAttributeVar).(string),
			AvatarUrlAttribute:         d.Get(idp_ldap.AvatarURLAttributeVar).(string),
			ProfileAttribute:           d.Get(idp_ldap.ProfileAttributeVar).(string),
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
	client, err := helper.GetManagementClient(clientinfo, d.Get(org_idp_utils.OrgIDVar).(string))
	if err != nil {
		return diag.FromErr(err)
	}
	timeout, err := time.ParseDuration(d.Get(idp_ldap.TimeoutVar).(string))
	if err != nil {
		return diag.FromErr(err)
	}
	_, err = client.UpdateLDAPProvider(ctx, &management.UpdateLDAPProviderRequest{
		Id:   d.Id(),
		Name: d.Get(idp_utils.NameVar).(string),
		ProviderOptions: &idp.Options{
			IsLinkingAllowed:  d.Get(idp_utils.IsLinkingAllowedVar).(bool),
			IsCreationAllowed: d.Get(idp_utils.IsCreationAllowedVar).(bool),
			IsAutoCreation:    d.Get(idp_utils.IsAutoCreationVar).(bool),
			IsAutoUpdate:      d.Get(idp_utils.IsAutoUpdateVar).(bool),
		},

		Servers:           idp_utils.InterfaceToStringSlice(d.Get(idp_ldap.ServersVar)),
		StartTls:          d.Get(idp_ldap.StartTLSVar).(bool),
		BaseDn:            d.Get(idp_ldap.BaseDNVar).(string),
		BindDn:            d.Get(idp_ldap.BindDNVar).(string),
		BindPassword:      d.Get(idp_ldap.BindPasswordVar).(string),
		UserBase:          d.Get(idp_ldap.UserBaseVar).(string),
		UserObjectClasses: helper.GetOkSetToStringSlice(d, idp_ldap.UserObjectClassesVar),
		UserFilters:       helper.GetOkSetToStringSlice(d, idp_ldap.UserFiltersVar),
		Timeout:           durationpb.New(timeout),

		Attributes: &idp.LDAPAttributes{
			IdAttribute:                d.Get(idp_ldap.IdAttributeVar).(string),
			FirstNameAttribute:         d.Get(idp_ldap.FirstNameAttributeVar).(string),
			LastNameAttribute:          d.Get(idp_ldap.LastNameAttributeVar).(string),
			DisplayNameAttribute:       d.Get(idp_ldap.DisplayNameAttributeVar).(string),
			NickNameAttribute:          d.Get(idp_ldap.NickNameAttributeVar).(string),
			PreferredUsernameAttribute: d.Get(idp_ldap.PreferredUsernameAttributeVar).(string),
			EmailAttribute:             d.Get(idp_ldap.EmailAttributeVar).(string),
			EmailVerifiedAttribute:     d.Get(idp_ldap.EmailVerifiedAttributeVar).(string),
			PhoneAttribute:             d.Get(idp_ldap.PhoneAttributeVar).(string),
			PhoneVerifiedAttribute:     d.Get(idp_ldap.PhoneVerifiedAttributeVar).(string),
			PreferredLanguageAttribute: d.Get(idp_ldap.PreferredLanguageAttributeVar).(string),
			AvatarUrlAttribute:         d.Get(idp_ldap.AvatarURLAttributeVar).(string),
			ProfileAttribute:           d.Get(idp_ldap.ProfileAttributeVar).(string),
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
	client, err := helper.GetManagementClient(clientinfo, d.Get(org_idp_utils.OrgIDVar).(string))
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
	idp := resp.GetIdp()
	cfg := idp.GetConfig()
	specificCfg := cfg.GetLdap()
	attributesCfg := specificCfg.GetAttributes()
	generalCfg := cfg.GetOptions()
	set := map[string]interface{}{
		org_idp_utils.OrgIDVar:         idp.GetDetails().GetResourceOwner(),
		idp_utils.NameVar:              idp.GetName(),
		idp_utils.IsLinkingAllowedVar:  generalCfg.GetIsLinkingAllowed(),
		idp_utils.IsCreationAllowedVar: generalCfg.GetIsCreationAllowed(),
		idp_utils.IsAutoCreationVar:    generalCfg.GetIsAutoCreation(),
		idp_utils.IsAutoUpdateVar:      generalCfg.GetIsAutoUpdate(),

		idp_ldap.ServersVar:           specificCfg.GetServers(),
		idp_ldap.StartTLSVar:          specificCfg.GetStartTls(),
		idp_ldap.BaseDNVar:            specificCfg.GetBaseDn(),
		idp_ldap.BindDNVar:            specificCfg.GetBindDn(),
		idp_ldap.BindPasswordVar:      d.Get(idp_ldap.BindPasswordVar).(string),
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
