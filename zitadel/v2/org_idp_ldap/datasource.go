package org_idp_ldap

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/helper"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/idp_ldap"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/idp_utils"
)

func GetDatasource() *schema.Resource {
	return &schema.Resource{
		Description: "Datasource representing an LDAP IdP on the organization.",
		Schema: map[string]*schema.Schema{
			idp_utils.IdpIDVar:             idp_utils.IdPIDDataSourceField,
			helper.OrgIDVar:                helper.OrgIDDatasourceField,
			idp_utils.NameVar:              idp_utils.NameDataSourceField,
			idp_utils.IsLinkingAllowedVar:  idp_utils.IsLinkingAllowedDataSourceField,
			idp_utils.IsCreationAllowedVar: idp_utils.IsCreationAllowedDataSourceField,
			idp_utils.IsAutoCreationVar:    idp_utils.IsAutoCreationDataSourceField,
			idp_utils.IsAutoUpdateVar:      idp_utils.IsAutoUpdateDataSourceField,

			idp_ldap.ServersVar:           idp_ldap.ServersDataSourceField,
			idp_ldap.StartTLSVar:          idp_ldap.StartTLSDataSourceField,
			idp_ldap.BaseDNVar:            idp_ldap.BaseDNDataSourceField,
			idp_ldap.BindDNVar:            idp_ldap.BindDNDataSourceField,
			idp_ldap.BindPasswordVar:      idp_ldap.BindPasswordDataSourceField,
			idp_ldap.UserBaseVar:          idp_ldap.UserBaseDataSourceField,
			idp_ldap.UserObjectClassesVar: idp_ldap.UserObjectClassesDataSourceField,
			idp_ldap.UserFiltersVar:       idp_ldap.UserFiltersDataSourceField,
			idp_ldap.TimeoutVar:           idp_ldap.TimeoutDataSourceField,
			idp_ldap.IdAttributeVar:       idp_ldap.IdAttributeDataSourceField,

			idp_ldap.FirstNameAttributeVar:         idp_ldap.FirstNameAttributeDataSourceField,
			idp_ldap.LastNameAttributeVar:          idp_ldap.LastNameAttributeDataSourceField,
			idp_ldap.DisplayNameAttributeVar:       idp_ldap.DisplayNameAttributeDataSourceField,
			idp_ldap.NickNameAttributeVar:          idp_ldap.NickNameAttributeDataSourceField,
			idp_ldap.PreferredUsernameAttributeVar: idp_ldap.PreferredUsernameAttributeDataSourceField,
			idp_ldap.EmailAttributeVar:             idp_ldap.EmailAttributeDataSourceField,
			idp_ldap.EmailVerifiedAttributeVar:     idp_ldap.EmailVerifiedAttributeDataSourceField,
			idp_ldap.PhoneAttributeVar:             idp_ldap.PhoneAttributeDataSourceField,
			idp_ldap.PhoneVerifiedAttributeVar:     idp_ldap.PhoneVerifiedAttributeDataSourceField,
			idp_ldap.PreferredLanguageAttributeVar: idp_ldap.PreferredLanguageAttributeDataSourceField,
			idp_ldap.AvatarURLAttributeVar:         idp_ldap.AvatarURLAttributeDataSourceField,
			idp_ldap.ProfileAttributeVar:           idp_ldap.ProfileAttributeDataSourceField,
		},
		ReadContext: read,
		Importer:    &schema.ResourceImporter{StateContext: helper.ImportWithIDAndOrgAndOptionalSecretStringV5(idp_utils.ClientSecretVar)},
	}
}
