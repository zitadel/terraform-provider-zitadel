package org_idp_github_es

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/idp_github_es"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/idp_utils"
)

func GetDatasource() *schema.Resource {
	return &schema.Resource{
		Description: "Datasource representing a GitHub Enterprise IdP of the organization.",
		Schema: map[string]*schema.Schema{
			idp_utils.IdpIDVar:                     idp_utils.IdPIDDataSourceField,
			helper.OrgIDVar:                        helper.OrgIDDatasourceField,
			idp_utils.NameVar:                      idp_utils.NameDataSourceField,
			idp_utils.ClientIDVar:                  idp_utils.ClientIDDataSourceField,
			idp_utils.ClientSecretVar:              idp_utils.ClientSecretDataSourceField,
			idp_utils.ScopesVar:                    idp_utils.ScopesDataSourceField,
			idp_utils.IsLinkingAllowedVar:          idp_utils.IsLinkingAllowedDataSourceField,
			idp_utils.IsCreationAllowedVar:         idp_utils.IsCreationAllowedDataSourceField,
			idp_utils.IsAutoCreationVar:            idp_utils.IsAutoCreationDataSourceField,
			idp_utils.IsAutoUpdateVar:              idp_utils.IsAutoUpdateDataSourceField,
			idp_utils.AutoLinkingVar:               idp_utils.AutoLinkingDataSourceField,
			idp_github_es.AuthorizationEndpointVar: idp_github_es.AuthorizationEndpointDataSourceField,
			idp_github_es.TokenEndpointVar:         idp_github_es.TokenEndpointDataSourceField,
			idp_github_es.UserEndpointVar:          idp_github_es.UserEndpointDataSourceField,
		},
		ReadContext: read,
	}
}
