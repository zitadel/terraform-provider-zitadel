package idp_ldap

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/helper"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/idp_utils"
)

func GetResource() *schema.Resource {
	return &schema.Resource{
		Description: "Resource representing an LDAP IDP on the instance.",
		Schema: map[string]*schema.Schema{
			idp_utils.NameVar:              idp_utils.NameResourceField,
			idp_utils.IsLinkingAllowedVar:  idp_utils.IsLinkingAllowedResourceField,
			idp_utils.IsCreationAllowedVar: idp_utils.IsCreationAllowedResourceField,
			idp_utils.IsAutoCreationVar:    idp_utils.IsAutoCreationResourceField,
			idp_utils.IsAutoUpdateVar:      idp_utils.IsAutoUpdateResourceField,

			ServersVar:           ServersResourceField,
			StartTLSVar:          StartTLSResourceField,
			BaseDNVar:            BaseDNResourceField,
			BindDNVar:            BindDNResourceField,
			BindPasswordVar:      BindPasswordResourceField,
			UserBaseVar:          UserBaseResourceField,
			UserObjectClassesVar: UserObjectClassesResourceField,
			UserFiltersVar:       UserFiltersResourceField,
			TimeoutVar:           TimeoutResourceField,
			IdAttributeVar:       IdAttributeResourceField,

			FirstNameAttributeVar:         FirstNameAttributeResourceField,
			LastNameAttributeVar:          LastNameAttributeResourceField,
			DisplayNameAttributeVar:       DisplayNameAttributeResourceField,
			NickNameAttributeVar:          NickNameAttributeResourceField,
			PreferredUsernameAttributeVar: PreferredUsernameAttributeResourceField,
			EmailAttributeVar:             EmailAttributeResourceField,
			EmailVerifiedAttributeVar:     EmailVerifiedAttributeResourceField,
			PhoneAttributeVar:             PhoneAttributeResourceField,
			PhoneVerifiedAttributeVar:     PhoneVerifiedAttributeResourceField,
			PreferredLanguageAttributeVar: PreferredLanguageAttributeResourceField,
			AvatarURLAttributeVar:         AvatarURLAttributeResourceField,
			ProfileAttributeVar:           ProfileAttributeResourceField,
		},
		ReadContext:   read,
		UpdateContext: update,
		CreateContext: create,
		DeleteContext: idp_utils.Delete,
		Importer:      helper.ImportWithIDAndOptionalSecret(idp_utils.IdpIDVar, BindPasswordVar),
	}
}
