package org_idp_saml

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/zitadel/terraform-provider-zitadel/zitadel/helper"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/idp_saml"
	"github.com/zitadel/terraform-provider-zitadel/zitadel/idp_utils"
)

func GetDatasource() *schema.Resource {
	return &schema.Resource{
		Description: "Datasource representing a SAML IdP of the organization.",
		Schema: map[string]*schema.Schema{
			idp_utils.IdpIDVar:             idp_utils.IdPIDDataSourceField,
			helper.OrgIDVar:                helper.OrgIDDatasourceField,
			idp_utils.NameVar:              idp_utils.NameDataSourceField,
			idp_saml.BindingVar:            idp_saml.BindingDatasourceField,
			idp_saml.MetadataXMLVar:        idp_saml.MetadataXMLDatasourceField,
			idp_saml.WithSignedRequestVar:  idp_saml.WithSignedRequestDatasourceField,
			idp_utils.IsLinkingAllowedVar:  idp_utils.IsLinkingAllowedDataSourceField,
			idp_utils.IsCreationAllowedVar: idp_utils.IsCreationAllowedDataSourceField,
			idp_utils.IsAutoCreationVar:    idp_utils.IsAutoCreationDataSourceField,
			idp_utils.IsAutoUpdateVar:      idp_utils.IsAutoUpdateDataSourceField,
		},
		ReadContext: read,
	}
}
