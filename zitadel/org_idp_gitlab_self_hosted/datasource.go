package org_idp_gitlab_self_hosted

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/helper"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/idp_gitlab_self_hosted"
	"github.com/zitadel/terraform-provider-zitadel/v2/zitadel/idp_utils"
)

func GetDatasource() *schema.Resource {
	return &schema.Resource{
		Description: "Datasource representing a GitLab Self Hosted IdP of the organization.",
		Schema: map[string]*schema.Schema{
			idp_utils.IdpIDVar:               idp_utils.IdPIDDataSourceField,
			helper.OrgIDVar:                  helper.OrgIDDatasourceField,
			idp_utils.NameVar:                idp_utils.NameDataSourceField,
			idp_utils.ClientIDVar:            idp_utils.ClientIDDataSourceField,
			idp_utils.ClientSecretVar:        idp_utils.ClientSecretDataSourceField,
			idp_utils.ScopesVar:              idp_utils.ScopesDataSourceField,
			idp_utils.IsLinkingAllowedVar:    idp_utils.IsLinkingAllowedDataSourceField,
			idp_utils.IsCreationAllowedVar:   idp_utils.IsCreationAllowedDataSourceField,
			idp_utils.IsAutoCreationVar:      idp_utils.IsAutoCreationDataSourceField,
			idp_utils.IsAutoUpdateVar:        idp_utils.IsAutoUpdateDataSourceField,
			idp_utils.AutoLinkingVar:         idp_utils.AutoLinkingDataSourceField,
			idp_gitlab_self_hosted.IssuerVar: idp_gitlab_self_hosted.IssuerDataSourceField,
		},
		ReadContext: read,
	}
}
