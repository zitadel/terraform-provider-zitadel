package org_idp_ldap

import (
	"context"
	"time"

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
	timeout, err := time.ParseDuration(d.Get(idp_utils.TimeoutVar).(string))
	if err != nil {
		return diag.FromErr(err)
	}
	resp, err := client.AddLDAPProvider(ctx, &management.AddLDAPProviderRequest{
		Name:              d.Get(idp_utils.NameVar).(string),
		Servers:           idp_utils.InterfaceToStringSlice(d.Get(idp_utils.ServersVar)),
		StartTls:          d.Get(idp_utils.StartTLSVar).(bool),
		BaseDn:            d.Get(idp_utils.BaseDNVar).(string),
		BindDn:            d.Get(idp_utils.BindDNVar).(string),
		BindPassword:      d.Get(idp_utils.BindPasswordVar).(string),
		UserBase:          d.Get(idp_utils.UserBaseVar).(string),
		UserObjectClasses: helper.GetOkSetToStringSlice(d, idp_utils.UserObjectClassesVar),
		UserFilters:       helper.GetOkSetToStringSlice(d, idp_utils.UserFiltersVar),
		Timeout:           durationpb.New(timeout),
		Attributes: &idp.LDAPAttributes{
			IdAttribute:                d.Get(idp_utils.IdAttributeVar).(string),
			FirstNameAttribute:         d.Get(idp_utils.FirstNameAttributeVar).(string),
			LastNameAttribute:          d.Get(idp_utils.LastNameAttributeVar).(string),
			DisplayNameAttribute:       d.Get(idp_utils.DisplayNameAttributeVar).(string),
			NickNameAttribute:          d.Get(idp_utils.NickNameAttributeVar).(string),
			PreferredUsernameAttribute: d.Get(idp_utils.PreferredUsernameAttributeVar).(string),
			EmailAttribute:             d.Get(idp_utils.EmailAttributeVar).(string),
			EmailVerifiedAttribute:     d.Get(idp_utils.EmailVerifiedAttributeVar).(string),
			PhoneAttribute:             d.Get(idp_utils.PhoneAttributeVar).(string),
			PhoneVerifiedAttribute:     d.Get(idp_utils.PhoneVerifiedAttributeVar).(string),
			PreferredLanguageAttribute: d.Get(idp_utils.PreferredLanguageAttributeVar).(string),
			AvatarUrlAttribute:         d.Get(idp_utils.AvatarURLAttributeVar).(string),
			ProfileAttribute:           d.Get(idp_utils.ProfileAttributeVar).(string),
		},
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
	client, err := helper.GetManagementClient(clientinfo, d.Get(org_idp_utils.OrgIDVar).(string))
	if err != nil {
		return diag.FromErr(err)
	}
	timeout, err := time.ParseDuration(d.Get(idp_utils.TimeoutVar).(string))
	if err != nil {
		return diag.FromErr(err)
	}
	if d.HasChangesExcept(idp_utils.IdpIDVar, org_idp_utils.OrgIDVar) {
		_, err = client.UpdateLDAPProvider(ctx, &management.UpdateLDAPProviderRequest{
			Id:                d.Id(),
			Name:              d.Get(idp_utils.NameVar).(string),
			Servers:           idp_utils.InterfaceToStringSlice(d.Get(idp_utils.ServersVar)),
			StartTls:          d.Get(idp_utils.StartTLSVar).(bool),
			BaseDn:            d.Get(idp_utils.BaseDNVar).(string),
			BindDn:            d.Get(idp_utils.BindDNVar).(string),
			BindPassword:      d.Get(idp_utils.BindPasswordVar).(string),
			UserBase:          d.Get(idp_utils.UserBaseVar).(string),
			UserObjectClasses: helper.GetOkSetToStringSlice(d, idp_utils.UserObjectClassesVar),
			UserFilters:       helper.GetOkSetToStringSlice(d, idp_utils.UserFiltersVar),
			Timeout:           durationpb.New(timeout),
			Attributes: &idp.LDAPAttributes{
				IdAttribute:                d.Get(idp_utils.IdAttributeVar).(string),
				FirstNameAttribute:         d.Get(idp_utils.FirstNameAttributeVar).(string),
				LastNameAttribute:          d.Get(idp_utils.LastNameAttributeVar).(string),
				DisplayNameAttribute:       d.Get(idp_utils.DisplayNameAttributeVar).(string),
				NickNameAttribute:          d.Get(idp_utils.NickNameAttributeVar).(string),
				PreferredUsernameAttribute: d.Get(idp_utils.PreferredUsernameAttributeVar).(string),
				EmailAttribute:             d.Get(idp_utils.EmailAttributeVar).(string),
				EmailVerifiedAttribute:     d.Get(idp_utils.EmailVerifiedAttributeVar).(string),
				PhoneAttribute:             d.Get(idp_utils.PhoneAttributeVar).(string),
				PhoneVerifiedAttribute:     d.Get(idp_utils.PhoneVerifiedAttributeVar).(string),
				PreferredLanguageAttribute: d.Get(idp_utils.PreferredLanguageAttributeVar).(string),
				AvatarUrlAttribute:         d.Get(idp_utils.AvatarURLAttributeVar).(string),
				ProfileAttribute:           d.Get(idp_utils.ProfileAttributeVar).(string),
			},
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
		org_idp_utils.OrgIDVar:                  idp.GetDetails().GetResourceOwner(),
		idp_utils.NameVar:                       idp.GetName(),
		idp_utils.IsLinkingAllowedVar:           generalCfg.GetIsLinkingAllowed(),
		idp_utils.IsCreationAllowedVar:          generalCfg.GetIsCreationAllowed(),
		idp_utils.IsAutoCreationVar:             generalCfg.GetIsAutoCreation(),
		idp_utils.IsAutoUpdateVar:               generalCfg.GetIsAutoUpdate(),
		idp_utils.ServersVar:                    specificCfg.GetServers(),
		idp_utils.StartTLSVar:                   specificCfg.GetStartTls(),
		idp_utils.BaseDNVar:                     specificCfg.GetBaseDn(),
		idp_utils.BindDNVar:                     specificCfg.GetBindDn(),
		idp_utils.BindPasswordVar:               d.Get(idp_utils.BindPasswordVar).(string),
		idp_utils.UserBaseVar:                   specificCfg.GetUserBase(),
		idp_utils.UserObjectClassesVar:          specificCfg.GetUserObjectClasses(),
		idp_utils.UserFiltersVar:                specificCfg.GetUserFilters(),
		idp_utils.TimeoutVar:                    specificCfg.GetTimeout().AsDuration().String(),
		idp_utils.IdAttributeVar:                attributesCfg.GetIdAttribute(),
		idp_utils.FirstNameAttributeVar:         attributesCfg.GetFirstNameAttribute(),
		idp_utils.LastNameAttributeVar:          attributesCfg.GetLastNameAttribute(),
		idp_utils.DisplayNameAttributeVar:       attributesCfg.GetDisplayNameAttribute(),
		idp_utils.NickNameAttributeVar:          attributesCfg.GetNickNameAttribute(),
		idp_utils.PreferredUsernameAttributeVar: attributesCfg.GetPreferredUsernameAttribute(),
		idp_utils.EmailAttributeVar:             attributesCfg.GetEmailAttribute(),
		idp_utils.EmailVerifiedAttributeVar:     attributesCfg.GetEmailVerifiedAttribute(),
		idp_utils.PhoneAttributeVar:             attributesCfg.GetPhoneAttribute(),
		idp_utils.PhoneVerifiedAttributeVar:     attributesCfg.GetPhoneVerifiedAttribute(),
		idp_utils.PreferredLanguageAttributeVar: attributesCfg.GetPreferredLanguageAttribute(),
		idp_utils.AvatarURLAttributeVar:         attributesCfg.GetAvatarUrlAttribute(),
		idp_utils.ProfileAttributeVar:           attributesCfg.GetProfileAttribute(),
	}
	for k, v := range set {
		if err := d.Set(k, v); err != nil {
			return diag.Errorf("failed to set %s of oidc idp: %v", k, err)
		}
	}
	d.SetId(idp.Id)
	return nil
}
