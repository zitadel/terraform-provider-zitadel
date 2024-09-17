package org_idp_oauth

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/idp_oauth"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/idp_utils"
)

func GetDatasource() *schema.Resource {
	return &schema.Resource{
		Description: "Datasource representing a generic OAuth2 IDP of the organization.",
		Schema: map[string]*schema.Schema{
			idp_utils.IdpIDVar:                 idp_utils.IdPIDDataSourceField,
			helper.OrgIDVar:                    helper.OrgIDDatasourceField,
			idp_utils.NameVar:                  idp_utils.NameDataSourceField,
			idp_utils.ClientIDVar:              idp_utils.ClientIDDataSourceField,
			idp_utils.ClientSecretVar:          idp_utils.ClientSecretDataSourceField,
			idp_oauth.AuthorizationEndpointVar: idp_oauth.AuthorizationEndpointDatasourceField,
			idp_oauth.TokenEndpointVar:         idp_oauth.TokenEndpointDatasourceField,
			idp_oauth.UserEndpointVar:          idp_oauth.UserEndpointDatasourceField,
			idp_oauth.IdAttributeVar:           idp_oauth.IdAttributeDatasourceField,
			idp_utils.ScopesVar:                idp_utils.ScopesDataSourceField,
			idp_utils.IsLinkingAllowedVar:      idp_utils.IsLinkingAllowedDataSourceField,
			idp_utils.IsCreationAllowedVar:     idp_utils.IsCreationAllowedDataSourceField,
			idp_utils.IsAutoCreationVar:        idp_utils.IsAutoCreationDataSourceField,
			idp_utils.IsAutoUpdateVar:          idp_utils.IsAutoUpdateDataSourceField,
		},
		ReadContext: read,
	}
}
