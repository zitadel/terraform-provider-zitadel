package idp_ldap

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/admin"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/idp"
	"google.golang.org/protobuf/types/known/durationpb"

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
	timeout, err := time.ParseDuration(idp_utils.StringValue(d, TimeoutVar))
	if err != nil {
		return diag.FromErr(err)
	}
	req := &admin.AddLDAPProviderRequest{
		Name:            idp_utils.StringValue(d, idp_utils.NameVar),
		ProviderOptions: idp_utils.ProviderOptionsValue(d),

		Servers:           idp_utils.InterfaceToStringSlice(d.Get(ServersVar)),
		StartTls:          idp_utils.BoolValue(d, StartTLSVar),
		BaseDn:            idp_utils.StringValue(d, BaseDNVar),
		BindDn:            idp_utils.StringValue(d, BindDNVar),
		BindPassword:      idp_utils.StringValue(d, BindPasswordVar),
		UserBase:          idp_utils.StringValue(d, UserBaseVar),
		UserObjectClasses: helper.GetOkSetToStringSlice(d, UserObjectClassesVar),
		UserFilters:       helper.GetOkSetToStringSlice(d, UserFiltersVar),
		Timeout:           durationpb.New(timeout),

		Attributes: &idp.LDAPAttributes{
			IdAttribute:                idp_utils.StringValue(d, IdAttributeVar),
			FirstNameAttribute:         idp_utils.StringValue(d, FirstNameAttributeVar),
			LastNameAttribute:          idp_utils.StringValue(d, LastNameAttributeVar),
			DisplayNameAttribute:       idp_utils.StringValue(d, DisplayNameAttributeVar),
			NickNameAttribute:          idp_utils.StringValue(d, NickNameAttributeVar),
			PreferredUsernameAttribute: idp_utils.StringValue(d, PreferredUsernameAttributeVar),
			EmailAttribute:             idp_utils.StringValue(d, EmailAttributeVar),
			EmailVerifiedAttribute:     idp_utils.StringValue(d, EmailVerifiedAttributeVar),
			PhoneAttribute:             idp_utils.StringValue(d, PhoneAttributeVar),
			PhoneVerifiedAttribute:     idp_utils.StringValue(d, PhoneVerifiedAttributeVar),
			PreferredLanguageAttribute: idp_utils.StringValue(d, PreferredLanguageAttributeVar),
			AvatarUrlAttribute:         idp_utils.StringValue(d, AvatarURLAttributeVar),
			ProfileAttribute:           idp_utils.StringValue(d, ProfileAttributeVar),
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
	client, err := helper.GetAdminClient(ctx, clientinfo)
	if err != nil {
		return diag.FromErr(err)
	}
	timeout, err := time.ParseDuration(idp_utils.StringValue(d, TimeoutVar))
	if err != nil {
		return diag.FromErr(err)
	}
	_, err = client.UpdateLDAPProvider(ctx, &admin.UpdateLDAPProviderRequest{
		Id:              d.Id(),
		Name:            idp_utils.StringValue(d, idp_utils.NameVar),
		ProviderOptions: idp_utils.ProviderOptionsValue(d),

		Servers:           idp_utils.InterfaceToStringSlice(d.Get(ServersVar)),
		StartTls:          idp_utils.BoolValue(d, StartTLSVar),
		BaseDn:            idp_utils.StringValue(d, BaseDNVar),
		BindDn:            idp_utils.StringValue(d, BindDNVar),
		BindPassword:      idp_utils.StringValue(d, BindPasswordVar),
		UserBase:          idp_utils.StringValue(d, UserBaseVar),
		UserObjectClasses: helper.GetOkSetToStringSlice(d, UserObjectClassesVar),
		UserFilters:       helper.GetOkSetToStringSlice(d, UserFiltersVar),
		Timeout:           durationpb.New(timeout),

		Attributes: &idp.LDAPAttributes{
			IdAttribute:                idp_utils.StringValue(d, IdAttributeVar),
			FirstNameAttribute:         idp_utils.StringValue(d, FirstNameAttributeVar),
			LastNameAttribute:          idp_utils.StringValue(d, LastNameAttributeVar),
			DisplayNameAttribute:       idp_utils.StringValue(d, DisplayNameAttributeVar),
			NickNameAttribute:          idp_utils.StringValue(d, NickNameAttributeVar),
			PreferredUsernameAttribute: idp_utils.StringValue(d, PreferredUsernameAttributeVar),
			EmailAttribute:             idp_utils.StringValue(d, EmailAttributeVar),
			EmailVerifiedAttribute:     idp_utils.StringValue(d, EmailVerifiedAttributeVar),
			PhoneAttribute:             idp_utils.StringValue(d, PhoneAttributeVar),
			PhoneVerifiedAttribute:     idp_utils.StringValue(d, PhoneVerifiedAttributeVar),
			PreferredLanguageAttribute: idp_utils.StringValue(d, PreferredLanguageAttributeVar),
			AvatarUrlAttribute:         idp_utils.StringValue(d, AvatarURLAttributeVar),
			ProfileAttribute:           idp_utils.StringValue(d, ProfileAttributeVar),
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
	specificCfg := cfg.GetLdap()
	attributesCfg := specificCfg.GetAttributes()
	generalCfg := cfg.GetOptions()
	set := map[string]interface{}{
		idp_utils.NameVar:              idp.GetName(),
		idp_utils.IsLinkingAllowedVar:  generalCfg.GetIsLinkingAllowed(),
		idp_utils.IsCreationAllowedVar: generalCfg.GetIsCreationAllowed(),
		idp_utils.IsAutoCreationVar:    generalCfg.GetIsAutoCreation(),
		idp_utils.IsAutoUpdateVar:      generalCfg.GetIsAutoUpdate(),
		idp_utils.AutoLinkingVar:       idp_utils.AutoLinkingString(generalCfg.GetAutoLinking()),

		ServersVar:           specificCfg.GetServers(),
		StartTLSVar:          specificCfg.GetStartTls(),
		BaseDNVar:            specificCfg.GetBaseDn(),
		BindDNVar:            specificCfg.GetBindDn(),
		BindPasswordVar:      idp_utils.StringValue(d, BindPasswordVar),
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
