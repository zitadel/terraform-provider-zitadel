package org_idp_azure_ad

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/helper"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/idp_azure_ad"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/v2/idp_utils"
)

func GetDatasource() *schema.Resource {
	return &schema.Resource{
		Description: "Datasource representing an Azure AD IdP of the organization.",
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
			idp_azure_ad.TenantTypeVar:     idp_azure_ad.TenantTypeDataSourceField,
			idp_azure_ad.TenantIDVar:       idp_azure_ad.TenantIDDataSourceField,
			idp_azure_ad.EmailVerifiedVar:  idp_azure_ad.EmailVerifiedDataSourceField,
		},
		ReadContext: read,
		Importer:    &schema.ResourceImporter{StateContext: helper.ImportWithIDAndOrgAndOptionalSecretStringV5(idp_utils.ClientSecretVar)},
	}
}
