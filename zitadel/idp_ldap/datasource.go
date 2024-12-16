package idp_ldap

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/idp_utils"
)

func GetDatasource() *schema.Resource {
	return &schema.Resource{
		Description: "Datasource representing an LDAP IDP on the instance.",
		Schema: map[string]*schema.Schema{
			idp_utils.IdpIDVar:             idp_utils.IdPIDDataSourceField,
			idp_utils.NameVar:              idp_utils.NameDataSourceField,
			idp_utils.IsLinkingAllowedVar:  idp_utils.IsLinkingAllowedDataSourceField,
			idp_utils.IsCreationAllowedVar: idp_utils.IsCreationAllowedDataSourceField,
			idp_utils.IsAutoCreationVar:    idp_utils.IsAutoCreationDataSourceField,
			idp_utils.IsAutoUpdateVar:      idp_utils.IsAutoUpdateDataSourceField,
			idp_utils.AutoLinkingVar:       idp_utils.AutoLinkingDataSourceField,

			ServersVar:           ServersDataSourceField,
			StartTLSVar:          StartTLSDataSourceField,
			BaseDNVar:            BaseDNDataSourceField,
			BindDNVar:            BindDNDataSourceField,
			BindPasswordVar:      BindPasswordDataSourceField,
			UserBaseVar:          UserBaseDataSourceField,
			UserObjectClassesVar: UserObjectClassesDataSourceField,
			UserFiltersVar:       UserFiltersDataSourceField,
			TimeoutVar:           TimeoutDataSourceField,
			IdAttributeVar:       IdAttributeDataSourceField,

			FirstNameAttributeVar:         FirstNameAttributeDataSourceField,
			LastNameAttributeVar:          LastNameAttributeDataSourceField,
			DisplayNameAttributeVar:       DisplayNameAttributeDataSourceField,
			NickNameAttributeVar:          NickNameAttributeDataSourceField,
			PreferredUsernameAttributeVar: PreferredUsernameAttributeDataSourceField,
			EmailAttributeVar:             EmailAttributeDataSourceField,
			EmailVerifiedAttributeVar:     EmailVerifiedAttributeDataSourceField,
			PhoneAttributeVar:             PhoneAttributeDataSourceField,
			PhoneVerifiedAttributeVar:     PhoneVerifiedAttributeDataSourceField,
			PreferredLanguageAttributeVar: PreferredLanguageAttributeDataSourceField,
			AvatarURLAttributeVar:         AvatarURLAttributeDataSourceField,
			ProfileAttributeVar:           ProfileAttributeDataSourceField,
		},
		ReadContext: read,
	}
}
