package org_idp_oidc

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/idp_oidc"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/idp_utils"
)

func GetDatasource() *schema.Resource {
	return &schema.Resource{
		Description: "Datasource representing a generic OIDC IdP on the organization.",
		Schema: map[string]*schema.Schema{
			idp_utils.IdpIDVar:             idp_utils.IdPIDDataSourceField,
			helper.OrgIDVar:                helper.OrgIDDatasourceField,
			idp_utils.NameVar:              idp_utils.NameDataSourceField,
			idp_utils.ClientIDVar:          idp_utils.ClientIDDataSourceField,
			idp_utils.ClientSecretVar:      idp_utils.ClientSecretDataSourceField,
			idp_utils.ScopesVar:            idp_utils.ScopesDataSourceField,
			idp_utils.IsLinkingAllowedVar:  idp_utils.IsLinkingAllowedDataSourceField,
			idp_utils.IsCreationAllowedVar: idp_utils.IsCreationAllowedDataSourceField,
			idp_utils.IsAutoCreationVar:    idp_utils.IsAutoCreationDataSourceField,
			idp_utils.IsAutoUpdateVar:      idp_utils.IsAutoUpdateDataSourceField,
			idp_utils.AutoLinkingVar:       idp_utils.AutoLinkingDataSourceField,
			idp_oidc.IssuerVar:             idp_oidc.IssuerDatasourceField,
			idp_oidc.IsIdTokenMappingVar:   idp_oidc.IsIdTokenMappingDatasourceField,
			idp_oidc.UsePKCEVar:            idp_oidc.UsePKCEDatasourceField,
		},
		ReadContext: read,
	}
}
