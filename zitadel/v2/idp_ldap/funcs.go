package idp_ldap

import (
	"context"
	"time"

	"google.golang.org/protobuf/types/known/durationpb"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/admin"
	"github.com/zitadel/zitadel-go/v2/pkg/client/zitadel/idp"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/helper"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/idp_utils"
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
	timeout, err := time.ParseDuration(d.Get(TimeoutVar).(string))
	if err != nil {
		return diag.FromErr(err)
	}
	req := &admin.AddLDAPProviderRequest{
		Name: d.Get(idp_utils.NameVar).(string),
		ProviderOptions: &idp.Options{
			IsLinkingAllowed:  d.Get(idp_utils.IsLinkingAllowedVar).(bool),
			IsCreationAllowed: d.Get(idp_utils.IsCreationAllowedVar).(bool),
			IsAutoUpdate:      d.Get(idp_utils.IsAutoUpdateVar).(bool),
			IsAutoCreation:    d.Get(idp_utils.IsAutoCreationVar).(bool),
		},

		Servers:           idp_utils.InterfaceToStringSlice(d.Get(ServersVar)),
		StartTls:          d.Get(StartTLSVar).(bool),
		BaseDn:            d.Get(BaseDNVar).(string),
		BindDn:            d.Get(BindDNVar).(string),
		BindPassword:      d.Get(BindPasswordVar).(string),
		UserBase:          d.Get(UserBaseVar).(string),
		UserObjectClasses: helper.GetOkSetToStringSlice(d, UserObjectClassesVar),
		UserFilters:       helper.GetOkSetToStringSlice(d, UserFiltersVar),
		Timeout:           durationpb.New(timeout),

		Attributes: &idp.LDAPAttributes{
			IdAttribute:                d.Get(IdAttributeVar).(string),
			FirstNameAttribute:         d.Get(FirstNameAttributeVar).(string),
			LastNameAttribute:          d.Get(LastNameAttributeVar).(string),
			DisplayNameAttribute:       d.Get(DisplayNameAttributeVar).(string),
			NickNameAttribute:          d.Get(NickNameAttributeVar).(string),
			PreferredUsernameAttribute: d.Get(PreferredUsernameAttributeVar).(string),
			EmailAttribute:             d.Get(EmailAttributeVar).(string),
			EmailVerifiedAttribute:     d.Get(EmailVerifiedAttributeVar).(string),
			PhoneAttribute:             d.Get(PhoneAttributeVar).(string),
			PhoneVerifiedAttribute:     d.Get(PhoneVerifiedAttributeVar).(string),
			PreferredLanguageAttribute: d.Get(PreferredLanguageAttributeVar).(string),
			AvatarUrlAttribute:         d.Get(AvatarURLAttributeVar).(string),
			ProfileAttribute:           d.Get(ProfileAttributeVar).(string),
		},
	}
	resp, err := client.AddLDAPProvider(ctx, req)
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
	timeout, err := time.ParseDuration(d.Get(TimeoutVar).(string))
	if err != nil {
		return diag.FromErr(err)
	}
	if d.HasChangesExcept(idp_utils.IdpIDVar) {
		_, err = client.UpdateLDAPProvider(ctx, &admin.UpdateLDAPProviderRequest{
			Id:   d.Id(),
			Name: d.Get(idp_utils.NameVar).(string),
			ProviderOptions: &idp.Options{
				IsLinkingAllowed:  d.Get(idp_utils.IsLinkingAllowedVar).(bool),
				IsCreationAllowed: d.Get(idp_utils.IsCreationAllowedVar).(bool),
				IsAutoCreation:    d.Get(idp_utils.IsAutoCreationVar).(bool),
				IsAutoUpdate:      d.Get(idp_utils.IsAutoUpdateVar).(bool),
			},

			Servers:           idp_utils.InterfaceToStringSlice(d.Get(ServersVar)),
			StartTls:          d.Get(StartTLSVar).(bool),
			BaseDn:            d.Get(BaseDNVar).(string),
			BindDn:            d.Get(BindDNVar).(string),
			BindPassword:      d.Get(BindPasswordVar).(string),
			UserBase:          d.Get(UserBaseVar).(string),
			UserObjectClasses: helper.GetOkSetToStringSlice(d, UserObjectClassesVar),
			UserFilters:       helper.GetOkSetToStringSlice(d, UserFiltersVar),
			Timeout:           durationpb.New(timeout),

			Attributes: &idp.LDAPAttributes{
				IdAttribute:                d.Get(IdAttributeVar).(string),
				FirstNameAttribute:         d.Get(FirstNameAttributeVar).(string),
				LastNameAttribute:          d.Get(LastNameAttributeVar).(string),
				DisplayNameAttribute:       d.Get(DisplayNameAttributeVar).(string),
				NickNameAttribute:          d.Get(NickNameAttributeVar).(string),
				PreferredUsernameAttribute: d.Get(PreferredUsernameAttributeVar).(string),
				EmailAttribute:             d.Get(EmailAttributeVar).(string),
				EmailVerifiedAttribute:     d.Get(EmailVerifiedAttributeVar).(string),
				PhoneAttribute:             d.Get(PhoneAttributeVar).(string),
				PhoneVerifiedAttribute:     d.Get(PhoneVerifiedAttributeVar).(string),
				PreferredLanguageAttribute: d.Get(PreferredLanguageAttributeVar).(string),
				AvatarUrlAttribute:         d.Get(AvatarURLAttributeVar).(string),
				ProfileAttribute:           d.Get(ProfileAttributeVar).(string),
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
	idp := resp.GetIdp()
	cfg := idp.GetConfig()
	specificCfg := cfg.GetLdap()
	attributesCfg := specificCfg.GetAttributes()
	generalCfg := cfg.GetOptions()
	set := map[string]interface{}{
		idp_utils.NameVar:              idp.GetName(),
		idp_utils.IsLinkingAllowedVar:  generalCfg.GetIsLinkingAllowed(),
		idp_utils.IsCreationAllowedVar: generalCfg.GetIsCreationAllowed(),
		idp_utils.IsAutoCreationVar:    generalCfg.GetIsAutoCreation(),
		idp_utils.IsAutoUpdateVar:      generalCfg.GetIsAutoUpdate(),

		ServersVar:           specificCfg.GetServers(),
		StartTLSVar:          specificCfg.GetStartTls(),
		BaseDNVar:            specificCfg.GetBaseDn(),
		BindDNVar:            specificCfg.GetBindDn(),
		BindPasswordVar:      d.Get(BindPasswordVar).(string),
		UserBaseVar:          specificCfg.GetUserBase(),
		UserObjectClassesVar: specificCfg.GetUserObjectClasses(),
		UserFiltersVar:       specificCfg.GetUserFilters(),
		TimeoutVar:           specificCfg.GetTimeout().AsDuration().String(),
		IdAttributeVar:       attributesCfg.GetIdAttribute(),

		FirstNameAttributeVar:         attributesCfg.GetFirstNameAttribute(),
		LastNameAttributeVar:          attributesCfg.GetLastNameAttribute(),
		DisplayNameAttributeVar:       attributesCfg.GetDisplayNameAttribute(),
		NickNameAttributeVar:          attributesCfg.GetNickNameAttribute(),
		PreferredUsernameAttributeVar: attributesCfg.GetPreferredUsernameAttribute(),
		EmailAttributeVar:             attributesCfg.GetEmailAttribute(),
		EmailVerifiedAttributeVar:     attributesCfg.GetEmailVerifiedAttribute(),
		PhoneAttributeVar:             attributesCfg.GetPhoneAttribute(),
		PhoneVerifiedAttributeVar:     attributesCfg.GetPhoneVerifiedAttribute(),
		PreferredLanguageAttributeVar: attributesCfg.GetPreferredLanguageAttribute(),
		AvatarURLAttributeVar:         attributesCfg.GetAvatarUrlAttribute(),
		ProfileAttributeVar:           attributesCfg.GetProfileAttribute(),
	}
	for k, v := range set {
		if err := d.Set(k, v); err != nil {
			return diag.Errorf("failed to set %s of oidc idp: %v", k, err)
		}
	}
	d.SetId(idp.Id)
	return nil
}
