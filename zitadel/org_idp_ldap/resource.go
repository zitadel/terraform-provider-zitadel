package org_idp_ldap

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/idp_ldap"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/idp_utils"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/org_idp_utils"
)

func GetResource() *schema.Resource {
	return &schema.Resource{
		Description: "Resource representing an LDAP IdP on the organization.",
		Schema: map[string]*schema.Schema{
			helper.OrgIDVar:                helper.OrgIDResourceField,
			idp_utils.NameVar:              idp_utils.NameResourceField,
			idp_utils.IsLinkingAllowedVar:  idp_utils.IsLinkingAllowedResourceField,
			idp_utils.IsCreationAllowedVar: idp_utils.IsCreationAllowedResourceField,
			idp_utils.IsAutoCreationVar:    idp_utils.IsAutoCreationResourceField,
			idp_utils.IsAutoUpdateVar:      idp_utils.IsAutoUpdateResourceField,
			idp_utils.AutoLinkingVar:       idp_utils.AutoLinkingResourceField,

			idp_ldap.ServersVar:           idp_ldap.ServersResourceField,
			idp_ldap.StartTLSVar:          idp_ldap.StartTLSResourceField,
			idp_ldap.BaseDNVar:            idp_ldap.BaseDNResourceField,
			idp_ldap.BindDNVar:            idp_ldap.BindDNResourceField,
			idp_ldap.BindPasswordVar:      idp_ldap.BindPasswordResourceField,
			idp_ldap.UserBaseVar:          idp_ldap.UserBaseResourceField,
			idp_ldap.UserObjectClassesVar: idp_ldap.UserObjectClassesResourceField,
			idp_ldap.UserFiltersVar:       idp_ldap.UserFiltersResourceField,
			idp_ldap.TimeoutVar:           idp_ldap.TimeoutResourceField,
			idp_ldap.IdAttributeVar:       idp_ldap.IdAttributeResourceField,
			idp_ldap.RootCAVar:            idp_ldap.RootCAResourceField,

			idp_ldap.FirstNameAttributeVar:         idp_ldap.FirstNameAttributeResourceField,
			idp_ldap.LastNameAttributeVar:          idp_ldap.LastNameAttributeResourceField,
			idp_ldap.DisplayNameAttributeVar:       idp_ldap.DisplayNameAttributeResourceField,
			idp_ldap.NickNameAttributeVar:          idp_ldap.NickNameAttributeResourceField,
			idp_ldap.PreferredUsernameAttributeVar: idp_ldap.PreferredUsernameAttributeResourceField,
			idp_ldap.EmailAttributeVar:             idp_ldap.EmailAttributeResourceField,
			idp_ldap.EmailVerifiedAttributeVar:     idp_ldap.EmailVerifiedAttributeResourceField,
			idp_ldap.PhoneAttributeVar:             idp_ldap.PhoneAttributeResourceField,
			idp_ldap.PhoneVerifiedAttributeVar:     idp_ldap.PhoneVerifiedAttributeResourceField,
			idp_ldap.PreferredLanguageAttributeVar: idp_ldap.PreferredLanguageAttributeResourceField,
			idp_ldap.AvatarURLAttributeVar:         idp_ldap.AvatarURLAttributeResourceField,
			idp_ldap.ProfileAttributeVar:           idp_ldap.ProfileAttributeResourceField,
		},
		ReadContext:   read,
		UpdateContext: update,
		CreateContext: create,
		DeleteContext: org_idp_utils.Delete,
		Importer:      helper.ImportWithIDAndOptionalOrgAndSecret(idp_utils.IdpIDVar, idp_ldap.BindPasswordVar),
	}
}
